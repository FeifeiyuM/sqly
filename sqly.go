package sqly

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// SqlY struct
type SqlY struct {
	db         *sql.DB
	driverName string
}

// Option sqly config option
type Option struct {
	Dsn             string        `json:"dsn"`                // database server name
	DriverName      string        `json:"driver_name"`        // database driver
	MaxIdleConns    int           `json:"max_idle_conns"`     // limit the number of idle connections
	MaxOpenConns    int           `json:"max_open_conns"`     // limit the number of total open connections
	ConnMaxLifeTime time.Duration `json:"conn_max_life_time"` // maximum amount of time a connection may be reused
}

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

// New init SqlY to database
func New(opt *Option) (*SqlY, error) {
	db, err := conn(opt.DriverName, opt.Dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(opt.ConnMaxLifeTime)
	db.SetMaxIdleConns(opt.MaxIdleConns)
	db.SetMaxOpenConns(opt.MaxOpenConns)
	return &SqlY{db: db, driverName: opt.DriverName}, nil
}

// Affected to record lastId for insert, and affected rows for update, inserts, delete statement
type Affected struct {
	LastId          int64 `json:"last_id"`
	LastIdErr       error `json:"last_id_Err"`
	RowsAffected    int64 `json:"rows_affected"`
	RowsAffectedErr error `json:"rows_affected_err"`
}

// exec one sql statement with context
func execOneDb(ctx context.Context, db *sql.DB, query string) (*Affected, error) {
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}
	aff := new(Affected)
	// last row_id that affected
	aff.LastId, aff.LastIdErr = res.LastInsertId()
	// rows that affected
	aff.RowsAffected, aff.RowsAffectedErr = res.RowsAffected()
	return aff, nil
}

// exec sql statements
func execManyDb(ctx context.Context, db *sql.DB, queries []string) error {
	// start transaction
	tx, err := db.Begin()
	defer func() {
		_ = tx.Rollback()
	}()
	if err != nil {
		return err
	}
	var errR error
	for _, query := range queries {
		_, errR = tx.ExecContext(ctx, query)
		if errR != nil {
			return errors.New("query:" + query + "; error:" + errR.Error())
		}
	}
	return tx.Commit()
}

// Ping ping test
func (s *SqlY) Ping() error {
	return s.db.Ping()
}

// Close close connection
func (s *SqlY) Close() error {
	return s.db.Close()
}

// Query query the database working with results
func (s *SqlY) Query(dest interface{}, query string, args ...interface{}) error {
	// query db
	q, err := queryFormat(query, args...)
	if err != nil {
		if errors.Is(err, ErrEmptyArrayInStatement) {
			return nil
		}
		return err
	}
	rows, err := s.db.Query(q)
	if err != nil {
		return err
	}
	return checkAllV2(rows, dest)
}

// Get query the database working with one result
func (s *SqlY) Get(dest interface{}, query string, args ...interface{}) error {
	// query db
	q, err := queryFormat(query, args...)
	if err != nil {
		if errors.Is(err, ErrEmptyArrayInStatement) {
			return nil
		}
		return err
	}
	rows, err := s.db.Query(q)
	if err != nil {
		return err
	}
	return checkOneV2(rows, dest)
}

// Insert insert into the database
func (s *SqlY) Insert(query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneDb(context.Background(), s.db, q)
}

// InsertMany insert many values to database
func (s *SqlY) InsertMany(query string, args [][]interface{}) (*Affected, error) {
	q, err := multiRowsFmt(query, args)
	if err != nil {
		return nil, err
	}
	return execOneDb(context.Background(), s.db, q)
}

// Update update value to database
func (s *SqlY) Update(query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneDb(context.Background(), s.db, q)
}

// UpdateMany update many
func (s *SqlY) UpdateMany(query string, args [][]interface{}) (*Affected, error) {
	var q string
	for _, arg := range args {
		t, err := queryFormat(query, arg...)
		if err != nil {
			return nil, err
		}
		q += t + ";"
	}
	return execOneDb(context.Background(), s.db, q)
}

// Delete delete item from database
func (s *SqlY) Delete(query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneDb(context.Background(), s.db, q)
}

// Exec general sql statement execute
func (s *SqlY) Exec(query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneDb(context.Background(), s.db, q)
}

// ExecMany execute multi sql statement
func (s *SqlY) ExecMany(queries []string) error {
	return execManyDb(context.Background(), s.db, queries)
}

// QueryCtx query the database working with results
func (s *SqlY) QueryCtx(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	// query db
	q, err := queryFormat(query, args...)
	if err != nil {
		if errors.Is(err, ErrEmptyArrayInStatement) {
			return nil
		}
		return err
	}
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	return checkAllV2(rows, dest)
}

// GetCtx query the database working with one result
func (s *SqlY) GetCtx(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	// query db
	q, err := queryFormat(query, args...)
	if err != nil {
		if errors.Is(err, ErrEmptyArrayInStatement) {
			return nil
		}
		return err
	}
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	return checkOneV2(rows, dest)
}

// InsertCtx insert with context
func (s *SqlY) InsertCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneDb(ctx, s.db, q)
}

// InsertManyCtx insert many with context
func (s *SqlY) InsertManyCtx(ctx context.Context, query string, args [][]interface{}) (*Affected, error) {
	q, err := multiRowsFmt(query, args)
	if err != nil {
		return nil, err
	}
	return execOneDb(ctx, s.db, q)
}

// UpdateCtx update with context
func (s *SqlY) UpdateCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneDb(ctx, s.db, q)
}

// UpdateManyCtx update many
func (s *SqlY) UpdateManyCtx(ctx context.Context, query string, args [][]interface{}) (*Affected, error) {
	var q string
	for _, arg := range args {
		t, err := queryFormat(query, arg...)
		if err != nil {
			return nil, err
		}
		q += t + ";"
	}
	return execOneDb(ctx, s.db, q)
}

// DeleteCtx delete with context
func (s *SqlY) DeleteCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneDb(ctx, s.db, q)
}

// ExecCtx general sql statement execute with context
func (s *SqlY) ExecCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	q, err := queryFormat(query, args...)
	if err != nil {
		return nil, err
	}
	return execOneDb(ctx, s.db, q)
}

// ExecManyCtx execute multi sql statement with context
func (s *SqlY) ExecManyCtx(ctx context.Context, queries []string) error {
	return execManyDb(ctx, s.db, queries)
}

// TxFunc callback function definition
type TxFunc func(tx *Trans) (interface{}, error)

// Transaction start transaction with callback function
func (s *SqlY) Transaction(txFunc TxFunc) (interface{}, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	// close or rollback transaction
	defer func() {
		_ = tx.Rollback()
	}()

	trans := Trans{tx}
	// run callback
	result, errR := txFunc(&trans)
	if errR != nil {
		return nil, errR
	}
	if errC := tx.Commit(); errC != nil {
		return nil, errC
	}
	return result, nil
}

// NewTrans start transaction
func (s *SqlY) NewTrans() (*Trans, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Trans{tx: tx}, nil
}
