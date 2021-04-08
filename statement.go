package sqly

import (
	"bytes"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func argToStr(delim string, item interface{}) (string, error) {
	if item == nil {
		return "NULL", nil
	}
	ref := reflect.Indirect(reflect.ValueOf(item)).Interface()
	switch v := ref.(type) {
	case bool:
		if v {
			return "1", nil
		} else {
			return "0", nil
		}
	case int:
		return strconv.Itoa(v), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case string:
		return SingleQuote(v), nil
	case time.Time:
		return SingleQuote(v.Format("2006-01-02 15:04:05.000000000")), nil
	case NullInt64:
		if v.Valid || v.Int64 != 0 {
			return strconv.FormatInt(v.Int64, 10), nil
		}
		return "NULL", nil
	case NullInt32:
		if v.Valid || v.Int32 != 0 {
			return strconv.FormatInt(int64(v.Int32), 10), nil
		}
		return "NULL", nil
	case NullFloat64:
		if v.Valid || v.Float64 != 0 {
			return strconv.FormatFloat(v.Float64, 'f', -1, 64), nil
		}
		return "NULL", nil
	case NullString:
		if v.Valid || v.String != "" {
			return SingleQuote(v.String), nil
		}
		return "NULL", nil
	case NullBool:
		if v.Valid {
			if v.Bool {
				return "1", nil
			}
			return "0", nil
		}
		return "NULL", nil
	case NullTime:
		if !v.Valid && v.Time.IsZero() {
			return "NULL", nil
		}
		return SingleQuote(v.Time.Format("2006-01-02 15:04:05.000000000")), nil
	case Boolean:
		if v {
			return "1", nil
		}
		return "0", nil
	case []int:
		if len(v) == 0 {
			return "", ErrEmptyArrayInStatement
		}
		var buffer bytes.Buffer
		buffer.WriteString("(")
		for i := 0; i < len(v); i++ {
			buffer.WriteString(strconv.Itoa(v[i]))
			if i != len(v)-1 {
				buffer.WriteString(delim)
			}
		}
		buffer.WriteString(")")
		return buffer.String(), nil
	case []int8:
		if len(v) == 0 {
			return "", ErrEmptyArrayInStatement
		}
		var buffer bytes.Buffer
		buffer.WriteString("(")
		for i := 0; i < len(v); i++ {
			buffer.WriteString(strconv.FormatInt(int64(v[i]), 10))
			if i != len(v)-1 {
				buffer.WriteString(delim)
			}
		}
		buffer.WriteString(")")
		return buffer.String(), nil
	case []int16:
		if len(v) == 0 {
			return "", ErrEmptyArrayInStatement
		}
		var buffer bytes.Buffer
		buffer.WriteString("(")
		for i := 0; i < len(v); i++ {
			buffer.WriteString(strconv.FormatInt(int64(v[i]), 10))
			if i != len(v)-1 {
				buffer.WriteString(delim)
			}
		}
		buffer.WriteString(")")
		return buffer.String(), nil
	case []int32:
		if len(v) == 0 {
			return "", ErrEmptyArrayInStatement
		}
		var buffer bytes.Buffer
		buffer.WriteString("(")
		for i := 0; i < len(v); i++ {
			buffer.WriteString(strconv.FormatInt(int64(v[i]), 10))
			if i != len(v)-1 {
				buffer.WriteString(delim)
			}
		}
		buffer.WriteString(")")
		return buffer.String(), nil
	case []int64:
		if len(v) == 0 {
			return "", ErrEmptyArrayInStatement
		}
		var buffer bytes.Buffer
		buffer.WriteString("(")
		for i := 0; i < len(v); i++ {
			buffer.WriteString(strconv.FormatInt(v[i], 10))
			if i != len(v)-1 {
				buffer.WriteString(delim)
			}
		}
		buffer.WriteString(")")
		return buffer.String(), nil
	case []float32:
		if len(v) == 0 {
			return "", ErrEmptyArrayInStatement
		}
		var buffer bytes.Buffer
		buffer.WriteString("(")
		for i := 0; i < len(v); i++ {
			buffer.WriteString(strconv.FormatFloat(float64(v[i]), 'f', -1, 32))
			if i != len(v)-1 {
				buffer.WriteString(delim)
			}
		}
		buffer.WriteString(")")
		return buffer.String(), nil
	case []float64:
		if len(v) == 0 {
			return "", ErrEmptyArrayInStatement
		}
		var buffer bytes.Buffer
		buffer.WriteString("(")
		for i := 0; i < len(v); i++ {
			buffer.WriteString(strconv.FormatFloat(v[i], 'f', -1, 64))
			if i != len(v)-1 {
				buffer.WriteString(delim)
			}
		}
		buffer.WriteString(")")
		return buffer.String(), nil
	case []string:
		if len(v) == 0 {
			return "", ErrEmptyArrayInStatement
		}
		var buffer bytes.Buffer
		buffer.WriteString("(")
		for i := 0; i < len(v); i++ {
			//buffer.WriteString("\"" + v[i] + "\"")
			buffer.WriteString(SingleQuote(v[i]))
			if i != len(v)-1 {
				buffer.WriteString(delim)
			}
		}
		buffer.WriteString(")")
		return buffer.String(), nil
	case []time.Time:
		if len(v) == 0 {
			return "", ErrEmptyArrayInStatement
		}
		var buffer bytes.Buffer
		buffer.WriteString("(")
		for i := 0; i < len(v); i++ {
			t := v[i].Format("2006-01-02 15:04:05.000000000")
			buffer.WriteString(SingleQuote(t))
			if i != len(v)-1 {
				buffer.WriteString(delim)
			}
		}
		buffer.WriteString(")")
		return buffer.String(), nil
	case []byte:
		return SingleQuote(string(v)), nil
	default:
		return "", ErrArgType
	}
}

// sql statement assemble
func queryFormat(fmtStr string, args ...interface{}) (string, error) {
	aLen := len(args)
	if aLen == 0 {
		return fmtStr, nil
	}
	fmtArr := strings.Split(fmtStr, "?")
	if len(fmtArr) == 1 {
		return fmtArr[0], nil
	}
	if len(fmtArr) != aLen+1 {
		return "", ErrQueryFmt
	}
	query := ""
	for idx, arg := range args {
		tmp, err := argToStr(",", arg)
		if err != nil {
			return "", err
		}
		query += fmtArr[idx] + tmp
	}
	query += fmtArr[aLen]
	return query, nil
}

// QueryFmt sql statement assemble public
func QueryFmt(fmtStr string, args ...interface{}) (string, error) {
	return queryFormat(fmtStr, args...)
}

// format rows that insert into a table
func multiRowsFmt(query string, args [][]interface{}) (string, error) {
	pat := `(\((\?,\s*)+\?*\s*\))`
	r, _ := regexp.Compile(pat)
	c := r.FindString(query)
	if c == "" {
		return "", ErrStatement
	}
	q := strings.Split(query, c)[0]

	var items []string
	for _, arg := range args {
		i, err := queryFormat(c, arg...)
		if err != nil {
			return "", err
		}
		items = append(items, i)
	}
	q += strings.Join(items, ",") + ";"
	return q, nil
}
