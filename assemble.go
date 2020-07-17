package sqly

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"
)

// close the Rows
func closeRows(rows *sql.Rows) error {
	if rows != nil {
		if err := rows.Close(); err != nil {
			return err
		}
	}
	return nil
}

var _scanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
var _timer = time.Time{}

// Indirect for reflect.Types
func directType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func baseType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	t = directType(t)
	if t.Kind() != expected {
		return nil, fmt.Errorf("expected %s but got %s", expected, t.Kind())
	}
	return t, nil
}

func scanAble(t reflect.Type) bool {
	if reflect.PtrTo(t).Implements(_scanner) {
		return true
	}
	if t == reflect.TypeOf(_timer) {
		return true
	}
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		return false
	}
	if t.Kind() == reflect.Struct || t.Kind() == reflect.Map {
		return false
	}
	return true
}

func isScanAble(field reflect.StructField) bool {
	return scanAble(field.Type)
}

func fieldsIterate(kvMap map[string][]int, pos []int, field reflect.StructField) {
	if !isScanAble(field) {
		var numField int
		var fieldType reflect.Type
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			numField = field.Type.Elem().NumField()
			fieldType = field.Type.Elem()
		} else {
			numField = field.Type.NumField()
			fieldType = field.Type
		}
		for i := 0; i < numField; i++ {
			_pos := append(pos, i)
			fieldsIterate(kvMap, _pos, fieldType.Field(i))
		}
	} else {
		tag := field.Tag.Get("sql")
		if tag == "" {
			tag = field.Name
		}
		kvMap[tag] = pos
	}
}

// fieldsMap
func fieldsColsMap(cols []string, mType reflect.Type) ([][]int, error) {
	kvMap := make(map[string][]int)
	for i := 0; i < mType.NumField(); i++ {
		pos := []int{i}
		fieldsIterate(kvMap, pos, mType.Field(i))
	}
	// span fields to list, to receive query values
	var fc [][]int
	for _, col := range cols {
		t, ok := kvMap[col]
		if !ok {
			//return nil, fmt.Errorf("field %s not exist", col)
			fc = append(fc, []int{-1})
		} else {
			fc = append(fc, t)
		}

	}
	return fc, nil
}

// fill values
func fieldAddrToContainer(v reflect.Value, fields [][]int, container []interface{}) error {
	//v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return errors.New("argument not a struct")
	}

	for i, pos := range fields {
		vt := v
		for si, p := range pos {
			// 处理接收字段少于数据库字段问题
			if p == -1 {
				container[i] = new(sql.RawBytes)
				continue
			}
			if vt.Kind() == reflect.Ptr && vt.Elem().Kind() == reflect.Struct {
				vt = vt.Elem().Field(p)
			} else {
				vt = vt.Field(p)
			}
			if vt.Kind() == reflect.Ptr && vt.IsNil() {
				alloc := reflect.New(directType(vt.Type()))
				vt.Set(alloc)
			}
			if si == len(pos)-1 {
				container[i] = vt.Addr().Interface()
				break
			}
		}
	}

	return nil
}

func allStructCheck(rows *sql.Rows, dVal reflect.Value, base reflect.Type, isPtr bool) error {
	// get columns name
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	// map column's name and container item fields
	fields, err := fieldsColsMap(cols, base)
	if err != nil {
		return err
	}

	// for store scan items
	con := make([]interface{}, len(cols))
	for rows.Next() {
		vp := reflect.New(base)
		v := reflect.Indirect(vp)

		err = fieldAddrToContainer(v, fields, con)
		if err != nil {
			return err
		}
		// scan val
		err = rows.Scan(con...)
		if err != nil {
			return err
		}
		if isPtr {
			dVal.Set(reflect.Append(dVal, vp))
		} else {
			dVal.Set(reflect.Append(dVal, v))
		}
	}
	return nil
}

func allBaseCheck(rows *sql.Rows, dVal reflect.Value, base reflect.Type, isPtr bool) error {
	for rows.Next() {
		vp := reflect.New(base)
		err := rows.Scan(vp.Interface())
		if err != nil {
			return err
		}
		// append
		if isPtr {
			dVal.Set(reflect.Append(dVal, vp))
		} else {
			dVal.Set(reflect.Append(dVal, reflect.Indirect(vp)))
		}
	}
	return nil
}

func allMapCheck(rows *sql.Rows, dVal reflect.Value) error {
	colsType, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	for rows.Next() {
		con := parseColumnsType(colsType)

		if err := rows.Scan(con...); err != nil {
			return err
		}
		v := make(map[string]interface{})
		for i, col := range colsType {
			v[col.Name()] = con[i].(interface{})
		}
		// append
		dVal.Set(reflect.Append(dVal, reflect.ValueOf(v)))
	}

	return nil
}

// scan all
func checkAllV2(rows *sql.Rows, dest interface{}) error {

	defer func() {
		if err := closeRows(rows); err != nil {
			panic(err)
		}
	}()

	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return ErrContainer
	}
	// construct container(dest) instance
	dVal := reflect.Indirect(val)

	// get container type instance
	dType, err := baseType(val.Type(), reflect.Slice)
	if err != nil {
		return err
	}
	// check container items is pointer type or not
	isPtr := dType.Elem().Kind() == reflect.Ptr
	// container item base type
	base := directType(dType.Elem())
	// on support struct of current
	if scanAble(base) {
		// base type （scan bale type)
		return allBaseCheck(rows, dVal, base, isPtr)
	} else if base.Kind() == reflect.Struct {
		// struct type
		return allStructCheck(rows, dVal, base, isPtr)
	} else if base.Kind() == reflect.Map {
		// map type
		return allMapCheck(rows, dVal)
	}
	return ErrContainer
}

func onlyCheck(st int) error {
	switch st {
	case 0:
		return ErrEmpty
	case 2:
		return ErrMultiRes
	default:
		return nil
	}
}

func structCheck(rows *sql.Rows, dVal reflect.Value, dType reflect.Type) error {
	// get columns name
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	// fields map
	fields, err := fieldsColsMap(cols, dType)
	if err != nil {
		return err
	}

	con := make([]interface{}, len(cols))
	err = fieldAddrToContainer(dVal, fields, con)
	if err != nil {
		return err
	}
	st := 0 // 0 no result, 1 one result, 2 more than one results
	for rows.Next() {
		st = 1
		if err := rows.Err(); err != nil {
			return err
		}
		if err := rows.Scan(con...); err != nil {
			return err
		}
		if rows.Next() {
			st = 2
			break
		}
	}
	return onlyCheck(st)
}

func mapCheck(rows *sql.Rows, dVal reflect.Value) error {
	colsType, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	columns := parseColumnsType(colsType)
	st := 0
	for rows.Next() {
		st = 1
		if err := rows.Err(); err != nil {
			return err
		}
		if err := rows.Scan(columns...); err != nil {
			return err
		}
		if rows.Next() {
			st = 2
			break
		}
	}
	if err := onlyCheck(st); err != nil {
		return err
	}
	for i, col := range colsType {
		dVal.SetMapIndex(reflect.ValueOf(col.Name()), reflect.ValueOf(columns[i]))
	}
	return nil
}

func baseTypeCheck(rows *sql.Rows, dest interface{}) error {
	st := 0
	for rows.Next() {
		st = 1
		err := rows.Scan(dest)
		if err != nil {
			return err
		}
		if rows.Next() {
			st = 2
			break
		}
	}
	return onlyCheck(st)
}

// query the database working with one result
func checkOneV2(rows *sql.Rows, dest interface{}) error {

	defer func() {
		if err := closeRows(rows); err != nil {
			panic(err)
		}
	}()

	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return ErrContainer
	}

	// construct container(dest) instance
	dVal := reflect.Indirect(val)
	if scanAble(dVal.Type()) {
		return baseTypeCheck(rows, dest)
	} else if dVal.Kind() == reflect.Struct {
		// struct type
		dType, err := baseType(val.Type(), reflect.Struct)
		if err != nil {
			return err
		}
		return structCheck(rows, dVal, dType)
	} else if dVal.Kind() == reflect.Map {
		// map type
		return mapCheck(rows, dVal)
	}
	return ErrContainer
}
