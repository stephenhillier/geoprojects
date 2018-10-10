package field

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
)

// Program represents one fieldwork job within a project.
// For example, a single trip that an engineer makes to collect field samples
// could be represented as a new Program record. There may be several new
// Datapoints (as well as Boreholes, Instruments) associated with that
// individual Program.
type Program struct {
	ID        int64     `json:"id"`
	Project   int64     `json:"project"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
}

// ProgramCreateRequest is the data a user should submit to create a program
type ProgramCreateRequest struct {
	Project   int64     `json:"project"`
	StartDate time.Time `json:"start_date" db:"start_date" schema:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date" schema:"end_date"`
}

// BoreholeCreateRequest is the data a user should submit to create a borehole.
// A borehole can either be associated with an existing datapoint, or if a location
// is supplied, a datapoint will be created.
type BoreholeCreateRequest struct {
	Project   int64      `json:"project"`
	Program   NullInt64  `json:"program"`
	Datapoint NullInt64  `json:"datapoint"`
	Name      string     `json:"name"`
	StartDate time.Time  `json:"start_date" db:"start_date" schema:"start_date"`
	EndDate   time.Time  `json:"end_date" db:"end_date" schema:"end_date"`
	FieldEng  NullInt64  `json:"field_eng" db:"field_eng" schema:"field_eng"`
	Location  [2]float64 `json:"location"`
}

// Datapoint is a geographic point that represents the location where data
// was gathered or where boreholes/instruments are located. A single datapoint
// may have a variety of data or records associated with it.
// Boreholes/instruments are organized this way (as children of a Datapoint) in
// order to reflect that drilling/subsurface sampling/instrumentation would
// generally have been performed at the same physical location.
type Datapoint struct {
	ID       NullInt64 `json:"id"`
	Location orb.Point `json:"location"`
}

// Borehole is drilled geotechnical test hole located at a Datapoint.
// There may be a number of samples/observations associated with one borehole.
type Borehole struct {
	ID        int64     `json:"id"`
	Project   NullInt64 `json:"project"`
	Program   NullInt64 `json:"program"`
	Datapoint NullInt64 `json:"datapoint"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
	FieldEng  NullInt64 `json:"field_eng" db:"field_eng"`
}

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON represents NullInt64 as JSON
func (v NullInt64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Int64)

}

// UnmarshalJSON converts from JSON to NullInt64
func (v *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &v.Int64)
	v.Valid = (err == nil)
	return err
}

// ProgramRepository is the set of methods that are available for interacting
// with field program records
type ProgramRepository interface {
	ListPrograms() ([]Program, error)
	CreateProgram(fp ProgramCreateRequest) (Program, error)
	GetProgram(programID int) (Program, error)
}

// DatapointRepository is the set of methods available for interacting with Datapoint records
type DatapointRepository interface {
	ListDatapoints() ([]*Datapoint, error)
	CreateDatapoint(dp Datapoint) (Datapoint, error)
	GetDatapoint(datapointID int) (Datapoint, error)
}

// BoreholeRepository is the set of methods available for interacting with Borehole records
type BoreholeRepository interface {
	ListBoreholes() ([]*Borehole, error)
	CreateBorehole(bh BoreholeCreateRequest) (Borehole, error)
	GetBorehole(boreholeID int) (Borehole, error)
}

// datastore is the database containing records related to the field module
// e.g. field programs, datapoints, boreholes
type datastore struct {
	*sqlx.DB
}

// Field Program database methods

// ListPrograms returns a list of field program records
func (db *datastore) ListPrograms() ([]Program, error) {
	programs := []Program{}
	query := `SELECT id, project, start_date, end_date FROM field_program`

	err := db.Select(&programs, query)
	if err != nil {
		return []Program{}, err
	}

	return programs, nil
}

// CreateProgram creates a field program record
func (db *datastore) CreateProgram(fp ProgramCreateRequest) (Program, error) {
	query := `INSERT INTO field_program (project, start_date, end_date) VALUES ($1, $2, $3) RETURNING id, project, start_date, end_date`
	created := Program{}
	err := db.Get(&created, query, fp.Project, fp.StartDate, fp.EndDate)
	if err != nil {
		return Program{}, err
	}
	return created, nil
}

// GetProgram retrieves a single field program record
func (db *datastore) GetProgram(programID int) (Program, error) {
	p := Program{}
	query := `SELECT id, project, start_date, end_date FROM field_program WHERE id=$1`
	err := db.Get(&p, query, programID)
	if err != nil {
		return Program{}, err
	}
	return p, nil
}

// Datapoint database methods

// CreateDatapoint creates a datapoint record.
// It may be called while handling create requests for boreholes or instruments
func (db *datastore) CreateDatapoint(dp Datapoint) (Datapoint, error) {
	query := `INSERT INTO datapoint (location) VALUES ($1) RETURNING id`
	created := Datapoint{}
	err := db.Get(&created, query, wkt.MarshalString(dp.Location))
	if err != nil {
		return Datapoint{}, err
	}
	return created, nil
}

// Borehole database methods

func (db *datastore) ListBoreholes() ([]*Borehole, error) {
	query := `SELECT id, project, program, start_date, end_date, field_eng FROM borehole`
	boreholes := []*Borehole{}
	err := db.Select(&boreholes, query)
	if err != nil {
		return []*Borehole{}, err
	}
	return boreholes, nil
}

// CreateBorehole creates a borehole record, as well as a Datapoint record if an existing
// datapoint wasn't supplied.
// Either a datapoint or a location should be supplied.
func (db *datastore) CreateBorehole(bh BoreholeCreateRequest) (Borehole, error) {

	// If a datapoint wasn't supplied, create one.
	// If a location also wasn't supplied, it will be created at the default location (0, 0?)
	if !bh.Datapoint.Valid {
		newDP := Datapoint{Location: orb.Point{bh.Location[0], bh.Location[1]}}
		createdDP, err := db.CreateDatapoint(newDP)
		if err != nil {
			return Borehole{}, err
		}
		bh.Datapoint = createdDP.ID
	}

	query := `INSERT INTO borehole (datapoint, program, project, name, start_date, end_date, field_eng) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, project, program, name, start_date, end_date, field_eng`
	created := Borehole{}
	err := db.Get(
		&created,
		query,
		bh.Datapoint,
		bh.Program,
		bh.Project,
		bh.Name,
		bh.StartDate,
		bh.EndDate,
		bh.FieldEng,
	)
	if err != nil {
		return Borehole{}, err
	}

	return created, nil
}

// GetBorehole retrieves a single borehole record.
func (db *datastore) GetBorehole(boreholeID int) (Borehole, error) {
	p := Borehole{}
	query := `SELECT id, project, program, datapoint, name, start_date, end_date, field_eng FROM borehole WHERE id=$1`
	err := db.Get(&p, query, boreholeID)
	if err != nil {
		return Borehole{}, err
	}
	return p, nil
}
