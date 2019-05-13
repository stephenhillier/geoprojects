package earthworks

import (
	"encoding/json"

	"github.com/stephenhillier/geoprojects/earthworks/db"
	"github.com/stephenhillier/geoprojects/earthworks/pkg/gis"
)

// BoreholeCreateRequest is the data a user should submit to create a borehole.
// A borehole can either be associated with an existing datapoint, or if a location
// is supplied, a datapoint will be created.
type BoreholeCreateRequest struct {
	Project        json.Number   `json:"project"`
	Program        db.NullInt64  `json:"program"`
	Datapoint      db.NullInt64  `json:"datapoint"`
	Name           string        `json:"name"`
	StartDate      db.NullDate   `json:"start_date" db:"start_date" schema:"start_date"`
	EndDate        db.NullDate   `json:"end_date" db:"end_date" schema:"end_date"`
	FieldEng       string        `json:"field_eng" db:"field_eng" schema:"field_eng"`
	Location       [2]float64    `json:"location"`
	DrillingMethod db.NullString `json:"drilling_method" db:"drilling_method"`
	TotalDepth     float64       `json:"total_depth" db:"total_depth"`
	Type           string        `json:"type" db:"type"` // type references borehole_type
}

// BoreholeResponse is the data returned by the API after receiving a request for
// a borehole's details
// the FieldEng field is a string (users.username) instead of a primary key reference.
type BoreholeResponse struct {
	ID                        int64             `json:"id"`
	Project                   db.NullInt64      `json:"project"`
	Program                   db.NullInt64      `json:"program"`
	Datapoint                 db.NullInt64      `json:"datapoint"`
	Name                      string            `json:"name"`
	StartDate                 db.NullDate       `json:"start_date" db:"start_date"`
	EndDate                   db.NullDate       `json:"end_date" db:"end_date"`
	FieldEng                  string            `json:"field_eng" db:"field_eng"`
	Location                  gis.PointLocation `json:"location"`
	StrataCount               int               `json:"strata_count,omitempty"`
	DrillingMethod            db.NullString     `json:"drilling_method" db:"drilling_method"`
	DrillingMethodDescription db.NullString     `json:"drilling_method_description" db:"drilling_method_description"`
	Type                      string            `json:"type" db:"type"` // type references borehole_type
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
	ID             int64         `json:"id"`
	Project        db.NullInt64  `json:"project"`
	Program        db.NullInt64  `json:"program"`
	Datapoint      int64         `json:"datapoint"`
	Name           string        `json:"name"`
	StartDate      db.NullDate   `json:"start_date" db:"start_date"`
	EndDate        db.NullDate   `json:"end_date" db:"end_date"`
	FieldEng       string        `json:"field_eng" db:"field_eng"`
	TotalDepth     float64       `json:"total_depth" db:"total_depth"`
	DrillingMethod db.NullString `json:"drilling_method" db:"drilling_method"`
	Type           string        `json:"type" db:"type"` // type references borehole_type
}

// Strata is a soil layer/stratum and contains information such as description and depth of the layer
type Strata struct {
	ID          int64   `json:"id"`
	Borehole    int64   `json:"borehole"`
	Start       float64 `json:"start" db:"start_depth"`
	End         float64 `json:"end" db:"end_depth"`
	Description string  `json:"description"`
	Soils       string  `json:"soils"`
	Moisture    string  `json:"moisture"`
	Consistency string  `json:"consistency"`
}

// StrataRequest is a struct containing fields required to create a new strata layer
type StrataRequest struct {
	Borehole    int64   `json:"borehole,string"`
	Start       float64 `json:"start,string"`
	End         float64 `json:"end,string"`
	Description string  `json:"description"`
}

// Sample is a soil layer/stratum and contains information such as description and depth of the layer
type Sample struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Borehole     int64   `json:"borehole"`
	BoreholeName string  `json:"borehole_name" db:"borehole_name"`
	Start        float64 `json:"start" db:"start_depth"`
	End          float64 `json:"end" db:"end_depth"`
	Description  string  `json:"description"`
	USCS         string  `json:"uscs" db:"uscs"`
	Tests        int     `json:"test_count"`
}

// SampleRequest is a struct containing fields required to create a new sample
type SampleRequest struct {
	Name        string  `json:"name"`
	Borehole    int64   `json:"borehole,string"`
	Start       float64 `json:"start,string"`
	End         float64 `json:"end,string"`
	Description string  `json:"description"`
	USCS        string  `json:"uscs" db:"uscs"`
}

// BoreholeRepository is the set of methods available for interacting with Borehole records
type BoreholeRepository interface {
	StrataRepository
	SampleRepository

	ListBoreholes(int, int, int) ([]*BoreholeResponse, int64, error)
	CreateBorehole(bh BoreholeCreateRequest, project int64) (Borehole, error)
	GetBorehole(boreholeID int) (BoreholeResponse, error)
	DeleteBorehole(boreholeID int64) error
}

// StrataRepository is the set of methods used for working with strata
// records in a database.
type StrataRepository interface {
	ListStrataByBorehole(boreholeID int64) ([]*Strata, error)
	CreateStrata(strata Strata) (Strata, error)
	CountStrataForBorehole(boreholeID int64) (int, error)
	RetrieveStrata(strataID int) (Strata, error)
	UpdateStrata(strata Strata) (Strata, error)
	DeleteStrata(strataID int64) error
}

// SampleRepository is the set of methods used for working with
// sample records in a database.
type SampleRepository interface {
	ListSamplesByBorehole(int64) ([]*Sample, error)
	ListSamplesByProject(int) ([]*Sample, error)
	CreateSample(Sample) (Sample, error)
	CountSampleForBorehole(int64) (int, error)
	RetrieveSample(int) (Sample, error)
	UpdateSample(Sample) (Sample, error)
	DeleteSample(int64) error
	CountTestForSample(int64) (int, error)
}

// BoreholeCtx is a context key for a borehole record
var BoreholeCtx struct{}

// StrataCtx is a context key for a strata record
var StrataCtx struct{}

// SampleCtx is a context key for a sample record in the request context
var SampleCtx struct{}
