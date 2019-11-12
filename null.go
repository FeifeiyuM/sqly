package sqlyt

import (
	"database/sql"
	"encoding/json"
)

// NullString is an alias for sql.NullTime
type NullTime sql.NullTime

// Scan implements the Scanner interface for NullTime
func (ns *NullTime) Scan(val interface{}) error {
	var t sql.NullTime
	if err := t.Scan(val); err != nil {
		return err
	}
	*ns = NullTime{t.Time, t.Valid}
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

// NullString is an alias for sql.NullBool
type NullBool sql.NullBool

// Scan implements the Scanner interface for NullBool
func (ns *NullBool) Scan(val interface{}) error {
	var b sql.NullBool
	if err := b.Scan(val); err != nil {
		return err
	}
	*ns = NullBool{b.Bool, b.Valid}
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

// NullString is an alias for sql.NullFloat64
type NullFloat64 sql.NullFloat64

// Scan implements the Scanner interface for NullFloat64
func (ns *NullFloat64) Scan(val interface{}) error {
	var f sql.NullFloat64
	if err := f.Scan(val); err != nil {
		return err
	}
	*ns = NullFloat64{f.Float64, f.Valid}
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

// NullString is an alias for sql.NullInt64
type NullInt64 sql.NullInt64

// Scan implements the Scanner interface for NullInt64
func (ns *NullInt64) Scan(val interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(val); err != nil {
		return err
	}
	*ns = NullInt64{i.Int64, i.Valid}
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

// NullString is an alias for sql.NullInt32
type NullInt32 sql.NullInt32

// Scan implements the Scanner interface for NullInt32
func (ns *NullInt32) Scan(val interface{}) error {
	var i sql.NullInt32
	if err := i.Scan(val); err != nil {
		return err
	}
	*ns = NullInt32{i.Int32, i.Valid}
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
	*ns = NullString{s.String, s.Valid}
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
