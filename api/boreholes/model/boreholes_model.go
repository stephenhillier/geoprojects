package model

import (
	"database/sql"
	"encoding/json"

	"github.com/stephenhillier/geoprojects/api/db"
	"github.com/stephenhillier/geoprojects/api/gis"
)

// BoreholeCreateRequest is the data a user should submit to create a borehole.
// A borehole can either be associated with an existing datapoint, or if a location
// is supplied, a datapoint will be created.
type BoreholeCreateRequest struct {
	Project   json.Number  `json:"project"`
	Program   db.NullInt64 `json:"program"`
	Datapoint db.NullInt64 `json:"datapoint"`
	Name      string       `json:"name"`
	StartDate db.NullDate  `json:"start_date" db:"start_date" schema:"start_date"`
	EndDate   db.NullDate  `json:"end_date" db:"end_date" schema:"end_date"`
	FieldEng  string       `json:"field_eng" db:"field_eng" schema:"field_eng"`
	Location  [2]float64   `json:"location"`
}

// BoreholeResponse is the data returned by the API after receiving a request for
// a borehole's details
// the FieldEng field is a string (users.username) instead of a primary key reference.
type BoreholeResponse struct {
	ID          int64             `json:"id"`
	Project     db.NullInt64      `json:"project"`
	Program     db.NullInt64      `json:"program"`
	Datapoint   db.NullInt64      `json:"datapoint"`
	Name        string            `json:"name"`
	StartDate   db.NullDate       `json:"start_date" db:"start_date"`
	EndDate     db.NullDate       `json:"end_date" db:"end_date"`
	FieldEng    string            `json:"field_eng" db:"field_eng"`
	Location    gis.PointLocation `json:"location"`
	StrataCount int               `json:"strata_count,omitempty"`
}

// Borehole is drilled geotechnical test hole located at a Datapoint.
// There may be a number of samples/observations associated with one borehole.
type Borehole struct {
	ID        int64        `json:"id"`
	Project   db.NullInt64 `json:"project"`
	Program   db.NullInt64 `json:"program"`
	Datapoint int64        `json:"datapoint"`
	Name      string       `json:"name"`
	StartDate db.NullDate  `json:"start_date" db:"start_date"`
	EndDate   db.NullDate  `json:"end_date" db:"end_date"`
	FieldEng  string       `json:"field_eng" db:"field_eng"`
}

// BoreholeV2 is a more comprehensive model that includes
// more data that might be required for professional use
type BoreholeV2 struct {
	ID             int64          `json:"id"`
	Project        db.NullInt64   `json:"project"`
	Program        db.NullInt64   `json:"program"`
	Datapoint      int64          `json:"datapoint"`
	Name           string         `json:"name"`
	StartDate      db.NullDate    `json:"start_date" db:"start_date"`
	EndDate        db.NullDate    `json:"end_date" db:"end_date"`
	FieldEng       string         `json:"field_eng" db:"field_eng"`
	TotalDepth     float64        `json:"total_depth" db:"total_depth"`
	DrillingMethod sql.NullString `json:"drilling_method" db:"drilling_method"`
}

// BoreholeCtx is a context key for a borehole record
var BoreholeCtx struct{}
