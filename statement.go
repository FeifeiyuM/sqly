package sqly

import (
	"bytes"
	"database/sql/driver"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func arrToStrPg(v driver.Value, err error) (string, error) {
	if err != nil {
		return "", err
	}
	if v == nil {
		return "NULL", nil
	}
	return PgString(v.(string)), nil
}

func arrToStrMySql(v driver.Value, err error) (string, error) {
	if err != nil {
		return "", err
	}
	if v == nil {
		return "NULL", nil
	}
	return SingleQuote(v.(string)), nil
}

type argFormat func(delim string, item interface{}) (string, error)

// mysql 参数 format
func mysqlArgFormat(delim string, item interface{}) (string, error) {
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
	case BoolArray:
		b, err := v.Value()
		return arrToStrMySql(b, err)
	case ByteaArray:
		b, err := v.Value()
		return arrToStrMySql(b, err)
	case Float64Array:
		b, err := v.Value()
		return arrToStrMySql(b, err)
	case Float32Array:
		b, err := v.Value()
		return arrToStrMySql(b, err)
	case GenericArray:
		b, err := v.Value()
		return arrToStrMySql(b, err)
	case Int64Array:
		b, err := v.Value()
		return arrToStrMySql(b, err)
	case StringArray:
		b, err := v.Value()
		return arrToStrMySql(b, err)
	default:
		return "", ErrArgType
	}
}

// postgresql 参数 format
func pgArgFormat(delim string, item interface{}) (string, error) {
	if item == nil {
		return "NULL", nil
	}
	ref := reflect.Indirect(reflect.ValueOf(item)).Interface()
	switch v := ref.(type) {
	case bool:
		if v {
			return "'t'", nil
		} else {
			return "'f'", nil
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
		return PgString(v), nil
	case time.Time:
		return PgString(v.Format("2006-01-02 15:04:05.000000000")), nil
	case NullInt64:
		if v.Valid {
			return strconv.FormatInt(v.Int64, 10), nil
		}
		return "NULL", nil
	case NullInt32:
		if v.Valid {
			return strconv.FormatInt(int64(v.Int32), 10), nil
		}
		return "NULL", nil
	case NullFloat64:
		if v.Valid {
			return strconv.FormatFloat(v.Float64, 'f', -1, 64), nil
		}
		return "NULL", nil
	case NullString:
		if v.Valid {
			return PgString(v.String), nil
		}
		return "NULL", nil
	case NullBool:
		if v.Valid {
			if v.Bool {
				return "'t'", nil
			}
			return "'f'", nil
		}
		return "NULL", nil
	case NullTime:
		if !v.Valid {
			return "NULL", nil
		}
		return PgString(v.Time.Format("2006-01-02 15:04:05.000000000")), nil
	case Boolean:
		if v {
			return "'t'", nil
		}
		return "'f'", nil
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
			buffer.WriteString(PgString(v[i]))
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
			buffer.WriteString(PgString(t))
			if i != len(v)-1 {
				buffer.WriteString(delim)
			}
		}
		buffer.WriteString(")")
		return buffer.String(), nil
	case []byte:
		return PgString(string(v)), nil
	case BoolArray:
		b, err := v.Value()
		return arrToStrPg(b, err)
	case ByteaArray:
		b, err := v.Value()
		return arrToStrPg(b, err)
	case Float64Array:
		b, err := v.Value()
		return arrToStrPg(b, err)
	case Float32Array:
		b, err := v.Value()
		return arrToStrPg(b, err)
	case GenericArray:
		b, err := v.Value()
		return arrToStrPg(b, err)
	case Int64Array:
		b, err := v.Value()
		return arrToStrPg(b, err)
	case StringArray:
		b, err := v.Value()
		return arrToStrPg(b, err)
	default:
		return "", ErrArgType
	}
}

// sql statement assemble
func statementFormat(fmtStr string, argFunc argFormat, args ...interface{}) (string, error) {
	if argFunc == nil {
		return fmtStr, nil
	}
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
		tmp, err := argFunc(",", arg)
		if err != nil {
			return "", err
		}
		query += fmtArr[idx] + tmp
	}
	query += fmtArr[aLen]
	return query, nil
}

// QueryFmtMysql sql statement assemble for mysql
func QueryFmtMysql(fmtStr string, args ...interface{}) (string, error) {
	return statementFormat(fmtStr, mysqlArgFormat, args...)
}

// format rows that insert into a table
func multiRowsFmt(query string, argFunc argFormat, args [][]interface{}) (string, error) {
	pat := `(\((\?,\s*)+\?*\s*\))`
	r, _ := regexp.Compile(pat)
	c := r.FindString(query)
	if c == "" {
		return "", ErrStatement
	}
	q := strings.Split(query, c)[0]

	var items []string
	for _, arg := range args {
		i, err := statementFormat(c, argFunc, arg...)
		if err != nil {
			return "", err
		}
		items = append(items, i)
	}
	q += strings.Join(items, ",") + ";"
	return q, nil
}

// QueryFmtPostgresql sql statement assemble for postgresql
func QueryFmtPostgresql(fmtStr string, args ...interface{}) (string, error) {
	return statementFormat(fmtStr, pgArgFormat, args...)
}
