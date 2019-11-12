package sqlyt

import (
	"context"
	"database/sql"
	"errors"
)

// sql struct for transaction
type Trans struct {
	tx *sql.Tx
}

// exec one sql statement with context
func execOneTx(ctx context.Context, tx *sql.Tx, query string) (*Affected, error) {
	res, err := tx.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}
	// last row_id that affected
	ld, errL := res.LastInsertId()
	// rows that affected
	rn, errR := res.RowsAffected()
	aff := new(Affected)
	if errL != nil && errR != nil {
		return nil, errors.New("LastInsertId err:" + errL.Error() + " RowsAffected err:" + errR.Error())
	}
	if errL != nil {
		aff.RowsAffected = rn
		return aff, errL
	}
	if errR != nil {
		aff.LastId = ld
		return nil, errR
	}
	aff.RowsAffected = rn
	aff.LastId = ld
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

// query results
func (t *Trans) Query(model interface{}, query string, args ...interface{}) (*[]interface{}, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	// query db
	rows, err := t.tx.Query(q)
	if err != nil {
		return nil, err
	}
	return checkAll(rows, model)
}

// query one row
func (t *Trans) QueryOne(model interface{}, query string, args ...interface{}) error {
	q, err := queryFormat(query, args...)
	if err != nil {
		return err
	}
	// query db
	rows, err := t.tx.Query(q)
	if err != nil {
		return err
	}
	return checkOne(rows, model)
}

// insert
func (t *Trans) Insert(query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneTx(context.Background(), t.tx, q)
}

// insert many rows
func (t *Trans) InsertMany(query string, args [][]interface{}) (*Affected, error) {
	q, err := multiRowsFmt(query, args)
	if err != nil {
		return nil, err
	}
	return execOneTx(context.Background(), t.tx, q)
}

// update
func (t *Trans) Update(query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneTx(context.Background(), t.tx, q)
}

// delete
func (t *Trans) Delete(query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneTx(context.Background(), t.tx, q)
}

// general sql statement execute
func (t *Trans) Exec(query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneTx(context.Background(), t.tx, q)
}

// execute multi sql statement
func (t *Trans) ExecMany(queries []string) error {
	return execManyTx(context.Background(), t.tx, queries)
}

// query results
func (t *Trans) QueryCtx(ctx context.Context, model interface{}, query string, args ...interface{}) (*[]interface{}, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	// query db
	rows, err := t.tx.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	return checkAll(rows, model)
}

// query one row
func (t *Trans) QueryOneCtx(ctx context.Context, model interface{}, query string, args ...interface{}) error {
	q, err := queryFormat(query, args...)
	if err != nil {
		return err
	}
	// query db
	rows, err := t.tx.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	return checkOne(rows, model)
}

// insert
func (t *Trans) InsertCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneTx(ctx, t.tx, q)
}

// insert many rows
func (t *Trans) InsertManyCtx(ctx context.Context, query string, args [][]interface{}) (*Affected, error) {
	q, err := multiRowsFmt(query, args)
	if err != nil {
		return nil, err
	}
	return execOneTx(ctx, t.tx, q)
}

// update
func (t *Trans) UpdateCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneTx(ctx, t.tx, q)
}

// delete
func (t *Trans) DeleteCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneTx(ctx, t.tx, q)
}

// general sql statement execute
func (t *Trans) ExecCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneTx(ctx, t.tx, q)
}

// execute multi sql statement
func (t *Trans) ExecManyCtx(ctx context.Context, queries []string) error {
	return execManyTx(ctx, t.tx, queries)
}
