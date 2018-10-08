package field

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/paulmach/orb"
)

// Program represents one fieldwork job within a project.
// For example, a single trip that an engineer makes to collect field samples
// could be represented as a new Program record. There may be several new
// Datapoints (as well as Boreholes, Instruments) associated with that
// individual Program.
type Program struct {
	ID        int       `json:"id"`
	Project   int       `json:"project"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
}

// ProgramCreateRequest is the data a user should submit to create a program
type ProgramCreateRequest struct {
	Project   int       `json:"project"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
}

// Datapoint is a geographic point that represents the location where data
// was gathered or where boreholes/instruments are located. A single datapoint
// may have a variety of data or records associated with it.
// Boreholes/instruments are organized this way (as children of a Datapoint) in
// order to reflect that drilling/subsurface sampling/instrumentation would
// generally have been performed at the same physical location.
type Datapoint struct {
	ID       int       `json:"id"`
	Location orb.Point `json:"location"`
}

// Borehole is drilled geotechnical test hole located at a Datapoint.
// There may be a number of samples/observations associated with one borehole.
type Borehole struct {
	ID        int       `json:"id"`
	Datapoint int       `json:"datapoint"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
	FieldEng  int       `json:"field_eng" db:"field_eng"`
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
	CreateBorehole(bh Borehole) (Borehole, error)
	GetBorehole(boreholeID int) (Borehole, error)
}

// datastore is the database containing records related to the field module
// e.g. field programs, datapoints, boreholes
type datastore struct {
	*sqlx.DB
}

func (db *datastore) ListPrograms() ([]Program, error) {
	programs := []Program{}
	query := `SELECT id, project, start_date, end_date FROM field_program`

	err := db.Select(&programs, query)
	if err != nil {
		return []Program{}, err
	}

	return programs, nil
}

func (db *datastore) CreateProgram(fp ProgramCreateRequest) (Program, error) {
	query := `INSERT INTO field_program (project, start_date, end_date) VALUES ($1, $2, $3) RETURNING id, project, start_date, end_date`
	created := Program{}
	err := db.Get(&created, query)
	if err != nil {
		return Program{}, err
	}
	return created, nil
}

func (db *datastore) GetProgram(programID int) (Program, error) {
	p := Program{}
	query := `SELECT id, project, start_date, end_date FROM field_program`
	err := db.Get(&p, query)
	if err != nil {
		return Program{}, err
	}
	return p, nil
}
