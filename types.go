package sqly

import (
	"database/sql"
	"encoding/json"
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
