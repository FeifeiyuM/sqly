package sqlyt

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
)

func argToStr(delim string, item interface{}) (string, error) {
	if item == nil {
		return "NULL", nil
	}
	ref := reflect.Indirect(reflect.ValueOf(item)).Interface()
	switch v := ref.(type) {
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
		return "\"" + v + "\"", nil
	case []int:
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
		var buffer bytes.Buffer
		buffer.WriteString("(")
		for i := 0; i < len(v); i++ {
			buffer.WriteString("\"" + v[i] + "\"")
			if i != len(v)-1 {
				buffer.WriteString(delim)
			}
		}
		buffer.WriteString(")")
		return buffer.String(), nil
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

// sql statement assemble public
func QueryFmt(fmtStr string, args ...interface{}) (string, error) {
	return queryFormat(fmtStr, args...)
}
