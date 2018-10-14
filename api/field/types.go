package field

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

// NullDate is an alias for pq.NullTime, and is meant to be used with dates only
type NullDate struct {
	pq.NullTime
}

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 struct {
	sql.NullInt64
}

// NullString is an alias for sql.NullString
type NullString struct {
	sql.NullString
}

// MarshalJSON represents NullInt64 as JSON
func (v NullInt64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Int64)
}

// MarshalJSON represents NullDate as JSON
func (v NullDate) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	layout := "2006-01-02"
	return json.Marshal(v.Time.Format(layout))
}

// UnmarshalJSON converts from JSON to NullDate
func (v *NullDate) UnmarshalJSON(b []byte) error {
	layout := "2006-01-02"
	var dateString string
	err := json.Unmarshal(b, &dateString)
	v.Valid = (err == nil)

	date, err := time.Parse(layout, dateString)
	v.Time = date
	v.Valid = (err == nil)

	return err
}

// MarshalJSON represents NullString as JSON
func (v NullString) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(v.String)
}

// UnmarshalJSON converts from JSON to NullInt64
func (v *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &v.Int64)
	v.Valid = (err == nil)
	return err
}

// UnmarshalJSON converts from JSON to NullString
func (v *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &v.String)
	v.Valid = (err == nil)
	return err
}
