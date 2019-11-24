package sqly

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
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

//// reflect retrieving result to struct
//func reflectModel(cols []string, mVal reflect.Value, mType reflect.Type) []interface{} {
//	kvMap := make(map[string]interface{})
//	for i := 0; i < mVal.NumField(); i++ {
//		vf := mVal.Field(i)
//		tf := mType.Field(i)
//		f := tf.Tag.Get("sql")
//		if f == "" {
//			f = tf.Name
//		}
//		kvMap[f] = vf.Addr().Interface()
//	}
//	// span fields to list, to receive query values
//	var scanDest []interface{}
//	for _, col := range cols {
//		scanDest = append(scanDest, kvMap[col])
//	}
//	return scanDest
//}
//
//// query the database working with results
//func checkAll(rows *sql.Rows, model interface{}) error {
//	modelType := reflect.TypeOf(model)
//	if modelType.Kind() != reflect.Ptr || modelType.Elem().Kind() != reflect.Slice {
//		return ErrContainer
//	}
//	modelVal := reflect.Indirect(reflect.ValueOf(model))
//	modelType = modelVal.Type()
//
//	cols, errC := rows.Columns()
//	if errC != nil {
//		return errC
//	}
//	// iterate over the rows
//	for rows.Next() {
//		if err := rows.Err(); err != nil {
//			return err
//		}
//		mType := modelType.Elem()
//		mVal := reflect.New(mType)
//		if mType.Kind() == reflect.Ptr {
//			fmt.Println("ok")
//			mVal = reflect.Indirect(mVal)
//		}
//		mVal = mVal.Elem()
//		scanDest := reflectModel(cols, mVal, mType)
//		// fanout results
//		if err := rows.Scan(scanDest...); err != nil {
//			return err
//		}
//		modelVal.Set(reflect.Append(modelVal, mVal))
//	}
//	// close rows
//	if err := closeRows(rows); err != nil {
//		return err
//	}
//	return nil
//}
//
//// query the database working with one result
//func checkOne(rows *sql.Rows, model interface{}) error {
//	modelType := reflect.TypeOf(model)
//	if modelType.Kind() != reflect.Ptr || modelType.Elem().Kind() != reflect.Struct {
//		return errors.New("invalid receive model")
//	}
//	modelVal := reflect.Indirect(reflect.ValueOf(model))
//	modelType = modelVal.Type()
//
//	cols, errC := rows.Columns()
//	if errC != nil {
//		return errC
//	}
//	for rows.Next() {
//		if err := rows.Err(); err != nil {
//			return err
//		}
//		scanDest := reflectModel(cols, modelVal, modelType)
//		if err := rows.Scan(scanDest...); err != nil {
//			return err
//		}
//		if rows.Next() {
//			break
//		}
//	}
//	// close rows
//	if err := closeRows(rows); err != nil {
//		return err
//	}
//	return nil
//}

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

// fieldsMap
func fieldsColsMap(cols []string, mType reflect.Type) []int {
	kvMap := make(map[string]int)
	for i := 0; i < mType.NumField(); i++ {
		tf := mType.Field(i)
		f := tf.Tag.Get("sql")
		if f == "" {
			f = tf.Name
		}
		kvMap[f] = i
	}
	// span fields to list, to receive query values
	var fc []int
	for _, col := range cols {
		fc = append(fc, kvMap[col])
	}
	return fc
}

// fill values
func fieldAddrToContainer(v reflect.Value, fields []int, container []interface{}) error {
	//v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return errors.New("argument not a struct")
	}

	for i, f := range fields {
		t := v.Field(f)
		// if this is a pointer and it's nil, allocate a new value and set it
		if t.Kind() == reflect.Ptr && t.IsNil() {
			alloc := reflect.New(directType(t.Type()))
			t.Set(alloc)
		}
		container[i] = t.Addr().Interface()
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
	// TODO to support int, string etc.
	if base.Kind() != reflect.Struct {
		return ErrContainer
	}
	// get columns name
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	// map column's name and container item fields
	fields := fieldsColsMap(cols, base)

	if len(cols) != len(fields) {
		return ErrFieldsMatch
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

	// get columns name
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	// get container type instance
	dType, err := baseType(val.Type(), reflect.Struct)
	if err != nil {
		return err
	}
	// fields map
	fields := fieldsColsMap(cols, dType)
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

	switch st {
	case 0:
		return ErrEmpty
	case 2:
		return ErrMultiRes
	default:
		return nil
	}
}
