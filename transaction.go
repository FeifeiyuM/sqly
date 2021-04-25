package sqly

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Trans sql struct for transaction
type Trans struct {
	tx     *sql.Tx
	driver dbDriver
}

// exec one sql statement with context
func (t *Trans) execOneTx(ctx context.Context, query string) (*Affected, error) {
	res, err := t.tx.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}
	aff := &Affected{
		result: res,
		driver: t.driver,
	}
	return aff, nil
}

// exec sql statements
func execManyTx(ctx context.Context, tx *sql.Tx, queries []string) error {
	for _, query := range queries {
		_, errR := tx.ExecContext(ctx, query)
		if errR != nil {
			return errors.New("query:" + query + "; error:" + errR.Error())
		}
	}
	return nil
}

// Rollback abort transaction
func (t *Trans) Rollback() error {
	return t.tx.Rollback()
}

// Commit commit transaction
func (t *Trans) Commit() error {
	return t.tx.Commit()
}

// Query query results
func (t *Trans) Query(dest interface{}, query string, args ...interface{}) error {
	q, err := statementFormat(query, argFmtFunc, args...)
	if err != nil {
		if errors.Is(err, ErrEmptyArrayInStatement) {
			return nil
		}
		return err
	}
	// query db
	rows, err := t.tx.Query(q)
	if err != nil {
		return err
	}
	return checkAllV2(rows, dest)
}

// Get query one row
func (t *Trans) Get(dest interface{}, query string, args ...interface{}) error {
	q, err := statementFormat(query, argFmtFunc, args...)
	if err != nil {
		if errors.Is(err, ErrEmptyArrayInStatement) {
			return nil
		}
		return err
	}
	// query db
	rows, err := t.tx.Query(q)
	if err != nil {
		return err
	}
	return checkOneV2(rows, dest)
}

// Insert insert
func (t *Trans) Insert(query string, args ...interface{}) (*Affected, error) {
	q, err := statementFormat(query, argFmtFunc, args...)
	if err != nil {
		return nil, err
	}
	return t.execOneTx(context.Background(), q)
}

// InsertMany insert many rows
func (t *Trans) InsertMany(query string, args [][]interface{}) (*Affected, error) {
	q, err := multiRowsFmt(query, argFmtFunc, args)
	if err != nil {
		return nil, err
	}
	return t.execOneTx(context.Background(), q)
}

// Update update
func (t *Trans) Update(query string, args ...interface{}) (*Affected, error) {
	q, err := statementFormat(query, argFmtFunc, args...)
	if err != nil {
		return nil, err
	}
	return t.execOneTx(context.Background(), q)
}

// UpdateMany update many
func (t *Trans) UpdateMany(query string, args [][]interface{}) (*Affected, error) {
	var qs []string
	for _, arg := range args {
		t, err := statementFormat(query, argFmtFunc, arg...)
		if err != nil {
			return nil, err
		}
		qs = append(qs, t)
	}
	q := strings.Join(qs, ";")
	return t.execOneTx(context.Background(), q)
}

// Delete delete
func (t *Trans) Delete(query string, args ...interface{}) (*Affected, error) {
	q, err := statementFormat(query, argFmtFunc, args...)
	if err != nil {
		return nil, err
	}
	return t.execOneTx(context.Background(), q)
}

// Exec general sql statement execute
func (t *Trans) Exec(query string, args ...interface{}) (*Affected, error) {
	q, err := statementFormat(query, argFmtFunc, args...)
	if err != nil {
		return nil, err
	}
	return t.execOneTx(context.Background(), q)
}

// ExecMany execute multi sql statement
func (t *Trans) ExecMany(queries []string) error {
	return execManyTx(context.Background(), t.tx, queries)
}

// QueryCtx query results
func (t *Trans) QueryCtx(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	q, err := statementFormat(query, argFmtFunc, args...)
	if err != nil {
		if errors.Is(err, ErrEmptyArrayInStatement) {
			return nil
		}
		return err
	}
	// query db
	rows, err := t.tx.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	return checkAllV2(rows, dest)
}

// GetCtx query one row
func (t *Trans) GetCtx(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	q, err := statementFormat(query, argFmtFunc, args...)
	if err != nil {
		if errors.Is(err, ErrEmptyArrayInStatement) {
			return nil
		}
		return err
	}
	// query db
	rows, err := t.tx.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	return checkOneV2(rows, dest)
}

// InsertCtx insert
func (t *Trans) InsertCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	q, err := statementFormat(query, argFmtFunc, args...)
	if err != nil {
		return nil, err
	}
	return t.execOneTx(ctx, q)
}

// InsertManyCtx insert many rows
func (t *Trans) InsertManyCtx(ctx context.Context, query string, args [][]interface{}) (*Affected, error) {
	q, err := multiRowsFmt(query, argFmtFunc, args)
	if err != nil {
		return nil, err
	}
	return t.execOneTx(ctx, q)
}

// UpdateCtx update
func (t *Trans) UpdateCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	q, err := statementFormat(query, argFmtFunc, args...)
	if err != nil {
		return nil, err
	}
	return t.execOneTx(ctx, q)
}

// UpdateManyCtx update many trans
func (t *Trans) UpdateManyCtx(ctx context.Context, query string, args [][]interface{}) (*Affected, error) {
	var q string
	for _, arg := range args {
		tmp, err := statementFormat(query, argFmtFunc, arg...)
		if err != nil {
			return nil, err
		}
		q += tmp + ";"
	}
	return t.execOneTx(ctx, q)
}

// DeleteCtx delete
func (t *Trans) DeleteCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	q, err := statementFormat(query, argFmtFunc, args...)
	if err != nil {
		return nil, err
	}
	return t.execOneTx(ctx, q)
}

// ExecCtx general sql statement execute
func (t *Trans) ExecCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	q, err := statementFormat(query, argFmtFunc, args...)
	if err != nil {
		return nil, err
	}
	return t.execOneTx(ctx, q)
}

// ExecManyCtx execute multi sql statement
func (t *Trans) ExecManyCtx(ctx context.Context, queries []string) error {
	return execManyTx(ctx, t.tx, queries)
}

// PgExec execute  statement for postgresql
func (t *Trans) PgExec(idField, query string, args ...interface{}) (*Affected, error) {
	return t.PgExecCtx(context.Background(), idField, query, args...)
}

// PgExecCtx execute  statement for postgresql with context
func (t *Trans) PgExecCtx(ctx context.Context, idField, query string, args ...interface{}) (*Affected, error) {
	if t.driver != driverPostgresql {
		return nil, ErrNotSupportForThisDriver
	}
	q, err := statementFormat(query, argFmtFunc, args...)
	if err != nil {
		return nil, err
	}
	q = fmt.Sprintf("%s RETURNING %s", q, idField)
	var id int64
	err = t.tx.QueryRowContext(ctx, q).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &Affected{
		lastId:       id,
		rowsAffected: -1,
		driver:       t.driver,
	}, nil
}
