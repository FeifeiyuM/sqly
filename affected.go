package sqly

import (
	"database/sql"
)

// Affected to record lastId for insert, and affected rows for update, inserts, delete statement
type Affected struct {
	result       sql.Result
	lastId       int64
	rowsAffected int64
	driverName   string
}

// GetLastId get lasted modified row id
func (a *Affected) GetLastId() (int64, error) {
	if a.lastId != 0 {
		return a.lastId, nil
	}
	if a.driverName == "postgres" {
		return 0, ErrNotSupportForThisDriver
	}
	var err error
	a.lastId, err = a.result.LastInsertId()
	return a.lastId, err
}

// GetRowsAffected returns the number of rows affected by an
// update, insert, or delete. Not every database or database
// driver may support this.
func (a *Affected) GetRowsAffected() (int64, error) {
	// local config of not support
	if a.rowsAffected == -1 {
		return 0, ErrNotSupportForThisDriver
	}
	if a.rowsAffected != 0 {
		return a.rowsAffected, nil
	}
	var err error
	a.rowsAffected, err = a.result.RowsAffected()
	return a.rowsAffected, err
}
