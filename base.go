package sqlyt

import (
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"strings"
)

// connect to database
func conn(driverName, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// close the Rows
func closeRows(rows *sql.Rows) error {
	if rows != nil {
		if err := rows.Close(); err != nil {
			return err
		}
	}
	return nil
}

// reflect retrieving result to struct
func reflectModel(cols []string, mVal reflect.Value, mType reflect.Type) []interface{} {
	kvMap := make(map[string]interface{})
	for i := 0; i < mVal.NumField(); i++ {
		vf := mVal.Field(i)
		tf := mType.Field(i)
		f := tf.Tag.Get("sql")
		if f == "" {
			f = tf.Name
		}
		kvMap[f] = vf.Addr().Interface()
	}
	// span fields to list, to receive query values
	var scanDest []interface{}
	for _, col := range cols {
		scanDest = append(scanDest, kvMap[col])
	}
	return scanDest
}

// query the database working with results
func checkAll(rows *sql.Rows, model interface{}) (*[]interface{}, error) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() != reflect.Ptr || modelType.Elem().Kind() != reflect.Struct {
		return nil, errors.New("invalid receive model")
	}

	cols, errC := rows.Columns()
	if errC != nil {
		return nil, errC
	}
	// result list
	var results []interface{}
	// iterate over the rows
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		mType := modelType.Elem()
		mVal := reflect.New(mType).Elem()
		scanDest := reflectModel(cols, mVal, mType)
		// fanout results
		if err := rows.Scan(scanDest...); err != nil {
			return nil, err
		}
		results = append(results, mVal.Interface())
	}
	// close rows
	if err := closeRows(rows); err != nil {
		return nil, err
	}
	return &results, nil
}

// query the database working with one result
func checkOne(rows *sql.Rows, model interface{}) error {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() != reflect.Ptr || modelType.Elem().Kind() != reflect.Struct {
		return errors.New("invalid receive model")
	}
	modelVal := reflect.Indirect(reflect.ValueOf(model))
	modelType = modelVal.Type()

	cols, errC := rows.Columns()
	if errC != nil {
		return errC
	}
	// iterate over the rows
	//if !rows.Next() {
	//	return sql.ErrNoRows
	//}
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		scanDest := reflectModel(cols, modelVal, modelType)
		if err := rows.Scan(scanDest...); err != nil {
			return err
		}
		break
	}
	// close rows
	if err := closeRows(rows); err != nil {
		return err
	}
	return nil
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

// Affected
// to record lastId for insert, and affected rows for update, inserts, delete statement
type Affected struct {
	LastId       int64
	RowsAffected int64
}

// errors
var ErrQueryFmt = errors.New("query can't be formatted")
var ErrArgType = errors.New("invalid variable type for argument")
var ErrStatement = errors.New("sql statement syntax error")
