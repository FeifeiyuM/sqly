/**
* origin from github.com/lib/pq array
* 并非原创
 */

package sqly

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// parseArray extracts the dimensions and elements of an array represented in
// text format. Only representations emitted by the backend are supported.
// Notably, whitespace around brackets and delimiters is significant, and NULL
// is case-sensitive.
//
// See http://www.postgresql.org/docs/current/static/arrays.html#ARRAYS-IO
func parseArray(src, del []byte) (dims []int, elems [][]byte, err error) {
	var depth, i int

	if len(src) < 1 || src[0] != '{' {
		return nil, nil, fmt.Errorf("pq: unable to parse array; expected %q at offset %d", '{', 0)
	}

Open:
	for i < len(src) {
		switch src[i] {
		case '{':
			depth++
			i++
		case '}':
			elems = make([][]byte, 0)
			goto Close
		default:
			break Open
		}
	}
	dims = make([]int, i)

Element:
	for i < len(src) {
		switch src[i] {
		case '{':
			if depth == len(dims) {
				break Element
			}
			depth++
			dims[depth-1] = 0
			i++
		case '"':
			var elem = []byte{}
			var escape bool
			for i++; i < len(src); i++ {
				if escape {
					elem = append(elem, src[i])
					escape = false
				} else {
					switch src[i] {
					default:
						elem = append(elem, src[i])
					case '\\':
						escape = true
					case '"':
						elems = append(elems, elem)
						i++
						break Element
					}
				}
			}
		default:
			for start := i; i < len(src); i++ {
				if bytes.HasPrefix(src[i:], del) || src[i] == '}' {
					elem := src[start:i]
					if len(elem) == 0 {
						return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
					}
					if bytes.Equal(elem, []byte("NULL")) {
						elem = nil
					}
					elems = append(elems, elem)
					break Element
				}
			}
		}
	}

	for i < len(src) {
		if bytes.HasPrefix(src[i:], del) && depth > 0 {
			dims[depth-1]++
			i += len(del)
			goto Element
		} else if src[i] == '}' && depth > 0 {
			dims[depth-1]++
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}

Close:
	for i < len(src) {
		if src[i] == '}' && depth > 0 {
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}
	if depth > 0 {
		err = fmt.Errorf("pq: unable to parse array; expected %q at offset %d", '}', i)
	}
	if err == nil {
		for _, d := range dims {
			if (len(elems) % d) != 0 {
				err = fmt.Errorf("pq: multidimensional arrays must have elements with matching dimensions")
			}
		}
	}
	return
}

func scanLinearArray(src, del []byte, typ string) (elems [][]byte, err error) {
	dims, elems, err := parseArray(src, del)
	if err != nil {
		return nil, err
	}
	if len(dims) > 1 {
		return nil, fmt.Errorf("pq: cannot convert ARRAY%s to %s", strings.Replace(fmt.Sprint(dims), " ", "][", -1), typ)
	}
	return elems, err
}

// Parse a bytea value received from the server.  Both "hex" and the legacy
// "escape" format are supported.
func parseBytea(s []byte) (result []byte, err error) {
	if len(s) >= 2 && bytes.Equal(s[:2], []byte("\\x")) {
		// bytea_output = hex
		s = s[2:] // trim off leading "\\x"
		result = make([]byte, hex.DecodedLen(len(s)))
		_, err := hex.Decode(result, s)
		if err != nil {
			return nil, err
		}
	} else {
		// bytea_output = escape
		for len(s) > 0 {
			if s[0] == '\\' {
				// escaped '\\'
				if len(s) >= 2 && s[1] == '\\' {
					result = append(result, '\\')
					s = s[2:]
					continue
				}

				// '\\' followed by an octal number
				if len(s) < 4 {
					return nil, fmt.Errorf("invalid bytea sequence %v", s)
				}
				r, err := strconv.ParseInt(string(s[1:4]), 8, 9)
				if err != nil {
					return nil, fmt.Errorf("could not parse bytea value: %s", err.Error())
				}
				result = append(result, byte(r))
				s = s[4:]
			} else {
				// We hit an unescaped, raw byte.  Try to read in as many as
				// possible in one go.
				i := bytes.IndexByte(s, '\\')
				if i == -1 {
					result = append(result, s...)
					break
				}
				result = append(result, s[:i]...)
				s = s[i:]
			}
		}
	}

	return result, nil
}

// appendArray appends rv to the buffer, returning the extended buffer and
// the delimiter used between elements.
//
// It panics when n <= 0 or rv's Kind is not reflect.Array nor reflect.Slice.
func appendArray(b []byte, rv reflect.Value, n int) ([]byte, string, error) {
	var del string
	var err error

	b = append(b, '{')

	if b, del, err = appendArrayElement(b, rv.Index(0)); err != nil {
		return b, del, err
	}

	for i := 1; i < n; i++ {
		b = append(b, del...)
		if b, del, err = appendArrayElement(b, rv.Index(i)); err != nil {
			return b, del, err
		}
	}

	return append(b, '}'), del, nil
}

// appendArrayElement appends rv to the buffer, returning the extended buffer
// and the delimiter to use before the next element.
//
// When rv's Kind is neither reflect.Array nor reflect.Slice, it is converted
// using driver.DefaultParameterConverter and the resulting []byte or string
// is double-quoted.
//
// See http://www.postgresql.org/docs/current/static/arrays.html#ARRAYS-IO

var typeByteSlice = reflect.TypeOf([]byte{})
var typeDriverValuer = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
var typeSQLScanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

// ArrayDelimiter may be optionally implemented by driver.Valuer or sql.Scanner
// to override the array delimiter used by GenericArray.
type ArrayDelimiter interface {
	// ArrayDelimiter returns the delimiter character(s) for this element's type.
	ArrayDelimiter() string
}

func appendArrayElement(b []byte, rv reflect.Value) ([]byte, string, error) {
	if k := rv.Kind(); k == reflect.Array || k == reflect.Slice {
		if t := rv.Type(); t != typeByteSlice && !t.Implements(typeDriverValuer) {
			if n := rv.Len(); n > 0 {
				return appendArray(b, rv, n)
			}

			return b, "", nil
		}
	}

	var del = ","
	var err error
	var iv interface{} = rv.Interface()

	if ad, ok := iv.(ArrayDelimiter); ok {
		del = ad.ArrayDelimiter()
	}

	if iv, err = driver.DefaultParameterConverter.ConvertValue(iv); err != nil {
		return b, del, err
	}

	switch v := iv.(type) {
	case nil:
		return append(b, "NULL"...), del, nil
	case []byte:
		return appendArrayQuotedBytes(b, v), del, nil
	case string:
		return appendArrayQuotedBytes(b, []byte(v)), del, nil
	}

	b, err = appendValue(b, iv)
	return b, del, err
}

func appendArrayQuotedBytes(b, v []byte) []byte {
	b = append(b, '"')
	for {
		i := bytes.IndexAny(v, `"\`)
		if i < 0 {
			b = append(b, v...)
			break
		}
		if i > 0 {
			b = append(b, v[:i]...)
		}
		b = append(b, '\\', v[i])
		v = v[i+1:]
	}
	return append(b, '"')
}

func appendValue(b []byte, v driver.Value) ([]byte, error) {
	return append(b, encode(v)...), nil
}

func encode(x interface{}) []byte {
	switch v := x.(type) {
	case int64:
		return strconv.AppendInt(nil, v, 10)
	case float64:
		return strconv.AppendFloat(nil, v, 'f', -1, 64)
	case []byte:
		return v
	case string:
		return []byte(v)
	case bool:
		return strconv.AppendBool(nil, v)
	case time.Time:
		return formatTs(v)

	default:
		panic(fmt.Errorf("encode: unknown type for %T", v))
	}
}

// formatTs formats t into a format postgres understands.
func formatTs(t time.Time) []byte {
	return FormatTimestamp(t)
}

// FormatTimestamp formats t into Postgres' text format for timestamps.
func FormatTimestamp(t time.Time) []byte {
	// Need to send dates before 0001 A.D. with " BC" suffix, instead of the
	// minus sign preferred by Go.
	// Beware, "0000" in ISO is "1 BC", "-0001" is "2 BC" and so on
	bc := false
	if t.Year() <= 0 {
		// flip year sign, and add 1, e.g: "0" will be "1", and "-10" will be "11"
		t = t.AddDate((-t.Year())*2+1, 0, 0)
		bc = true
	}
	b := []byte(t.Format("2006-01-02 15:04:05.999999999Z07:00"))

	_, offset := t.Zone()
	offset %= 60
	if offset != 0 {
		// RFC3339Nano already printed the minus sign
		if offset < 0 {
			offset = -offset
		}

		b = append(b, ':')
		if offset < 10 {
			b = append(b, '0')
		}
		b = strconv.AppendInt(b, int64(offset), 10)
	}

	if bc {
		b = append(b, " BC"...)
	}
	return b
}

// Array returns the optimal driver.Valuer and sql.Scanner for an array or
// slice of any dimension.
//
// For example:
//  db.Query(`SELECT * FROM t WHERE id = ANY($1)`, pq.Array([]int{235, 401}))
//
//  var x []sql.NullInt64
//  db.QueryRow(`SELECT ARRAY[235, 401]`).Scan(pq.Array(&x))
//
// Scanning multi-dimensional arrays is not supported.  Arrays where the lower
// bound is not one (such as `[0:0]={1}') are not supported.
func Array(a interface{}) interface {
	driver.Valuer
	sql.Scanner
} {
	switch a := a.(type) {
	case []bool:
		return (*BoolArray)(&a)
	case []float64:
		return (*Float64Array)(&a)
	case []float32:
		return (*Float32Array)(&a)
	case []int64:
		return (*Int64Array)(&a)
	case []int32:
		return (*Int32Array)(&a)
	case []string:
		return (*StringArray)(&a)
	case [][]byte:
		return (*ByteaArray)(&a)

	case *[]bool:
		return (*BoolArray)(a)
	case *[]float64:
		return (*Float64Array)(a)
	case *[]float32:
		return (*Float32Array)(a)
	case *[]int64:
		return (*Int64Array)(a)
	case *[]int32:
		return (*Int32Array)(a)
	case *[]string:
		return (*StringArray)(a)
	case *[][]byte:
		return (*ByteaArray)(a)
	}

	return GenericArray{a}
}
