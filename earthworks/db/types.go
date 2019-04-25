package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/lib/pq"
)

// Context Keys
type contextKey struct{}

var projectCtx contextKey
var boreholeCtx contextKey
var strataCtx contextKey
var sampleCtx contextKey
var labTestCtx contextKey

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

// NullFloat64 is an alias for sql.NullFloat64
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON represents NullInt64 as JSON
func (v NullInt64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Int64)
}

// MarshalJSON represents NullFloat64 as JSON
func (v NullFloat64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Float64)
}

// MarshalJSON represents NullDate as JSON
func (v NullDate) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	layout := "2006-01-02T15:04:05Z07:00"
	return json.Marshal(v.Time.Format(layout))
}

// UnmarshalJSON converts from JSON to NullDate
func (v *NullDate) UnmarshalJSON(b []byte) error {
	layout := "2006-01-02T15:04:05Z07:00"
	var dateString string
	err := json.Unmarshal(b, &dateString)
	v.Valid = (err == nil)

	date, err := time.Parse(layout, dateString)
	v.Time = date
	v.Valid = (err == nil)

	// If we didn't get a valid date using the full layout, try again with
	// a simplified layout. (YYYY-MM-DD format)
	if !v.Valid {
		simpleLayout := "2006-01-02"
		date, err = time.Parse(simpleLayout, dateString)
		v.Time = date
		v.Valid = (err == nil)
	}

	log.Println(date, err)
	return nil
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

// UnmarshalJSON converts from JSON to NullFloat64
func (v *NullFloat64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &v.Float64)
	v.Valid = (err == nil)
	return err
}

// UnmarshalJSON converts from JSON to NullString
func (v *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &v.String)
	v.Valid = (err == nil)
	return err
}
