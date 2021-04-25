package sqly

import (
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// NullTime is an alias for sql.NullTime
type NullTime sql.NullTime

// Scan implements the Scanner interface for NullTime
func (ns *NullTime) Scan(val interface{}) error {
	var t sql.NullTime
	if err := t.Scan(val); err != nil {
		return err
	}
	*ns = NullTime{Time: t.Time, Valid: t.Valid}
	return nil
}

// MarshalJSON for NullTime
func (ns *NullTime) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Time)
}

// UnmarshalJSON for NullTime
func (ns *NullTime) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.Time)
	ns.Valid = err == nil
	return err
}

// NullBool is an alias for sql.NullBool
type NullBool sql.NullBool

// Scan implements the Scanner interface for NullBool
func (ns *NullBool) Scan(val interface{}) error {
	var b sql.NullBool
	if err := b.Scan(val); err != nil {
		return err
	}
	*ns = NullBool{Bool: b.Bool, Valid: b.Valid}
	return nil
}

// MarshalJSON for NullBool
func (ns *NullBool) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Bool)
}

// UnmarshalJSON for NullBool
func (ns *NullBool) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.Bool)
	ns.Valid = err == nil
	return err
}

// NullFloat64 is an alias for sql.NullFloat64
type NullFloat64 sql.NullFloat64

// Scan implements the Scanner interface for NullFloat64
func (ns *NullFloat64) Scan(val interface{}) error {
	var f sql.NullFloat64
	if err := f.Scan(val); err != nil {
		return err
	}
	*ns = NullFloat64{Float64: f.Float64, Valid: f.Valid}
	return nil
}

// MarshalJSON for NullFloat64
func (ns *NullFloat64) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Float64)
}

// UnmarshalJSON for NullFloat64
func (ns *NullFloat64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.Float64)
	ns.Valid = err == nil
	return err
}

// NullInt64 is an alias for sql.NullInt64
type NullInt64 sql.NullInt64

// Scan implements the Scanner interface for NullInt64
func (ns *NullInt64) Scan(val interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(val); err != nil {
		return err
	}
	*ns = NullInt64{Int64: i.Int64, Valid: i.Valid}
	return nil
}

// MarshalJSON for NullInt64
func (ns *NullInt64) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Int64)
}

// UnmarshalJSON for NullInt64
func (ns *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.Int64)
	ns.Valid = err == nil
	return err
}

// NullInt32 is an alias for sql.NullInt32
type NullInt32 sql.NullInt32

// Scan implements the Scanner interface for NullInt32
func (ns *NullInt32) Scan(val interface{}) error {
	var i sql.NullInt32
	if err := i.Scan(val); err != nil {
		return err
	}
	*ns = NullInt32{Int32: i.Int32, Valid: i.Valid}
	return nil
}

// MarshalJSON for NullInt32
func (ns *NullInt32) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Int32)
}

// UnmarshalJSON for NullInt32
func (ns *NullInt32) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.Int32)
	ns.Valid = err == nil
	return err
}

// NullString is an alias for sql.NullString
type NullString sql.NullString

// Scan implements the Scanner interface for NullString
func (ns *NullString) Scan(val interface{}) error {
	var s sql.NullString
	if err := s.Scan(val); err != nil {
		return err
	}
	*ns = NullString{String: s.String, Valid: s.Valid}
	return nil
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = err == nil
	return err
}

// Boolean boolean
type Boolean bool

// Scan boolean Scan
func (b *Boolean) Scan(val interface{}) error {
	if val == nil || string(val.([]byte)) == "0" {
		*b = false
	} else {
		*b = true
	}
	return nil
}

// array types
// origin from github.com/lib/pq array

// BoolArray represents a one-dimensional array of the PostgreSQL boolean type.
type BoolArray []bool

// Scan implements the sql.Scanner interface.
func (a *BoolArray) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to BoolArray", src)
}

func (a *BoolArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "BoolArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(BoolArray, len(elems))
		for i, v := range elems {
			if len(v) != 1 {
				return fmt.Errorf("pq: could not parse boolean array index %d: invalid boolean %q", i, v)
			}
			switch v[0] {
			case 't':
				b[i] = true
			case 'f':
				b[i] = false
			default:
				return fmt.Errorf("pq: could not parse boolean array index %d: invalid boolean %q", i, v)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a BoolArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be exactly two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1+2*n)

		for i := 0; i < n; i++ {
			b[2*i] = ','
			if a[i] {
				b[1+2*i] = 't'
			} else {
				b[1+2*i] = 'f'
			}
		}

		b[0] = '{'
		b[2*n] = '}'

		return string(b), nil
	}

	return "{}", nil
}

// ByteaArray represents a one-dimensional array of the PostgreSQL bytea type.
type ByteaArray [][]byte

// Scan implements the sql.Scanner interface.
func (a *ByteaArray) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to ByteaArray", src)
}

func (a *ByteaArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "ByteaArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(ByteaArray, len(elems))
		for i, v := range elems {
			b[i], err = parseBytea(v)
			if err != nil {
				return fmt.Errorf("could not parse bytea array index %d: %s", i, err.Error())
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface. It uses the "hex" format which
// is only supported on PostgreSQL 9.0 or newer.
func (a ByteaArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, 2*N bytes of quotes,
		// 3*N bytes of hex formatting, and N-1 bytes of delimiters.
		size := 1 + 6*n
		for _, x := range a {
			size += hex.EncodedLen(len(x))
		}

		b := make([]byte, size)

		for i, s := 0, b; i < n; i++ {
			o := copy(s, `,"\\x`)
			o += hex.Encode(s[o:], a[i])
			s[o] = '"'
			s = s[o+1:]
		}

		b[0] = '{'
		b[size-1] = '}'

		return string(b), nil
	}

	return "{}", nil
}

// Float64Array represents a one-dimensional array of the PostgreSQL double
// precision type.
type Float64Array []float64

// Scan implements the sql.Scanner interface.
func (a *Float64Array) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to Float64Array", src)
}

func (a *Float64Array) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "Float64Array")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(Float64Array, len(elems))
		for i, v := range elems {
			if b[i], err = strconv.ParseFloat(string(v), 64); err != nil {
				return fmt.Errorf("pq: parsing array element index %d: %v", i, err)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a Float64Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendFloat(b, a[0], 'f', -1, 64)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendFloat(b, a[i], 'f', -1, 64)
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}

// Float32Array represents a one-dimensional array of the PostgreSQL double
// precision type.
type Float32Array []float32

// Scan implements the sql.Scanner interface.
func (a *Float32Array) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to Float32Array", src)
}

func (a *Float32Array) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "Float32Array")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(Float32Array, len(elems))
		for i, v := range elems {
			var x float64
			if x, err = strconv.ParseFloat(string(v), 32); err != nil {
				return fmt.Errorf("pq: parsing array element index %d: %v", i, err)
			}
			b[i] = float32(x)
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a Float32Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendFloat(b, float64(a[0]), 'f', -1, 32)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendFloat(b, float64(a[i]), 'f', -1, 32)
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}

// GenericArray implements the driver.Valuer and sql.Scanner interfaces for
// an array or slice of any dimension.
type GenericArray struct{ A interface{} }

func (GenericArray) evaluateDestination(rt reflect.Type) (reflect.Type, func([]byte, reflect.Value) error, string) {
	var assign func([]byte, reflect.Value) error
	var del = ","

	// TODO calculate the assign function for other types
	// TODO repeat this section on the element type of arrays or slices (multidimensional)
	{
		if reflect.PtrTo(rt).Implements(typeSQLScanner) {
			// dest is always addressable because it is an element of a slice.
			assign = func(src []byte, dest reflect.Value) (err error) {
				ss := dest.Addr().Interface().(sql.Scanner)
				if src == nil {
					err = ss.Scan(nil)
				} else {
					err = ss.Scan(src)
				}
				return
			}
			goto FoundType
		}

		assign = func([]byte, reflect.Value) error {
			return fmt.Errorf("pq: scanning to %s is not implemented; only sql.Scanner", rt)
		}
	}

FoundType:

	if ad, ok := reflect.Zero(rt).Interface().(ArrayDelimiter); ok {
		del = ad.ArrayDelimiter()
	}

	return rt, assign, del
}

// Scan implements the sql.Scanner interface.
func (a GenericArray) Scan(src interface{}) error {
	dpv := reflect.ValueOf(a.A)
	switch {
	case dpv.Kind() != reflect.Ptr:
		return fmt.Errorf("pq: destination %T is not a pointer to array or slice", a.A)
	case dpv.IsNil():
		return fmt.Errorf("pq: destination %T is nil", a.A)
	}

	dv := dpv.Elem()
	switch dv.Kind() {
	case reflect.Slice:
	case reflect.Array:
	default:
		return fmt.Errorf("pq: destination %T is not a pointer to array or slice", a.A)
	}

	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src, dv)
	case string:
		return a.scanBytes([]byte(src), dv)
	case nil:
		if dv.Kind() == reflect.Slice {
			dv.Set(reflect.Zero(dv.Type()))
			return nil
		}
	}

	return fmt.Errorf("pq: cannot convert %T to %s", src, dv.Type())
}

func (a GenericArray) scanBytes(src []byte, dv reflect.Value) error {
	dtype, assign, del := a.evaluateDestination(dv.Type().Elem())
	dims, elems, err := parseArray(src, []byte(del))
	if err != nil {
		return err
	}

	// TODO allow multidimensional

	if len(dims) > 1 {
		return fmt.Errorf("pq: scanning from multidimensional ARRAY%s is not implemented",
			strings.Replace(fmt.Sprint(dims), " ", "][", -1))
	}

	// Treat a zero-dimensional array like an array with a single dimension of zero.
	if len(dims) == 0 {
		dims = append(dims, 0)
	}

	for i, rt := 0, dv.Type(); i < len(dims); i, rt = i+1, rt.Elem() {
		switch rt.Kind() {
		case reflect.Slice:
		case reflect.Array:
			if rt.Len() != dims[i] {
				return fmt.Errorf("pq: cannot convert ARRAY%s to %s",
					strings.Replace(fmt.Sprint(dims), " ", "][", -1), dv.Type())
			}
		default:
			// TODO handle multidimensional
		}
	}

	values := reflect.MakeSlice(reflect.SliceOf(dtype), len(elems), len(elems))
	for i, e := range elems {
		if err := assign(e, values.Index(i)); err != nil {
			return fmt.Errorf("pq: parsing array element index %d: %v", i, err)
		}
	}

	// TODO handle multidimensional

	switch dv.Kind() {
	case reflect.Slice:
		dv.Set(values.Slice(0, dims[0]))
	case reflect.Array:
		for i := 0; i < dims[0]; i++ {
			dv.Index(i).Set(values.Index(i))
		}
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (a GenericArray) Value() (driver.Value, error) {
	if a.A == nil {
		return nil, nil
	}

	rv := reflect.ValueOf(a.A)

	switch rv.Kind() {
	case reflect.Slice:
		if rv.IsNil() {
			return nil, nil
		}
	case reflect.Array:
	default:
		return nil, fmt.Errorf("pq: Unable to convert %T to array", a.A)
	}

	if n := rv.Len(); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 0, 1+2*n)

		b, _, err := appendArray(b, rv, n)
		return string(b), err
	}

	return "{}", nil
}

// Int64Array represents a one-dimensional array of the PostgreSQL integer types.
type Int64Array []int64

// Scan implements the sql.Scanner interface.
func (a *Int64Array) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to Int64Array", src)
}

func (a *Int64Array) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "Int64Array")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(Int64Array, len(elems))
		for i, v := range elems {
			if b[i], err = strconv.ParseInt(string(v), 10, 64); err != nil {
				return fmt.Errorf("pq: parsing array element index %d: %v", i, err)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a Int64Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendInt(b, a[0], 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendInt(b, a[i], 10)
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}

// Int32Array represents a one-dimensional array of the PostgreSQL integer types.
type Int32Array []int32

// Scan implements the sql.Scanner interface.
func (a *Int32Array) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to Int32Array", src)
}

func (a *Int32Array) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "Int32Array")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(Int32Array, len(elems))
		for i, v := range elems {
			var x int
			if x, err = strconv.Atoi(string(v)); err != nil {
				return fmt.Errorf("pq: parsing array element index %d: %v", i, err)
			}
			b[i] = int32(x)
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a Int32Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendInt(b, int64(a[0]), 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendInt(b, int64(a[i]), 10)
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}

// StringArray represents a one-dimensional array of the PostgreSQL character types.
type StringArray []string

// Scan implements the sql.Scanner interface.
func (a *StringArray) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to StringArray", src)
}

func (a *StringArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "StringArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(StringArray, len(elems))
		for i, v := range elems {
			if b[i] = string(v); v == nil {
				return fmt.Errorf("pq: parsing array element index %d: cannot convert nil to string", i)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, 2*N bytes of quotes,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+3*n)
		b[0] = '{'

		b = appendArrayQuotedBytes(b, []byte(a[0]))
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = appendArrayQuotedBytes(b, []byte(a[i]))
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}

// ColumnsType Working with Unknown Columns
// mysql only
func parseColumnsType(colsType []*sql.ColumnType) []interface{} {
	var cTypes []interface{}

	for _, ct := range colsType {
		nullAble, _ := ct.Nullable()
		switch strings.ToUpper(ct.DatabaseTypeName()) {
		case "MEDIUMINT", "INT", "INTEGER", "BIGINT":
			if nullAble {
				cTypes = append(cTypes, new(NullInt64))
			} else {
				cTypes = append(cTypes, new(int64))
			}
		case "SMALLINT", "TINYINT":
			if nullAble {
				cTypes = append(cTypes, new(NullInt32))
			} else {
				cTypes = append(cTypes, new(int32))
			}
		case "FLOAT", "DOUBLE", "DECIMAL":
			if nullAble {
				cTypes = append(cTypes, new(NullFloat64))
			} else {
				cTypes = append(cTypes, new(float64))
			}
		case "DATE", "TIME", "YEAR", "DATETIME", "TIMESTAMP":
			if nullAble {
				cTypes = append(cTypes, new(NullTime))
			} else {
				cTypes = append(cTypes, new(time.Time))
			}
		case "CHAR", "VARCHAR", "TINYTEXT", "TEXT", "MEDIUMTEXT", "LONGTEXT":
			if nullAble {
				cTypes = append(cTypes, new(NullString))
			} else {
				cTypes = append(cTypes, new(string))
			}
		case "TINYBLOB", "BLOB", "MEDIUMBLOB", "LONGBLOB":
			cTypes = append(cTypes, new(sql.RawBytes))
		default:
			cTypes = append(cTypes, new(sql.RawBytes))
		}
	}
	return cTypes
}
