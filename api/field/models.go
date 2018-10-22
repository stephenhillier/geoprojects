package field

import (
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
	ID        int64    `json:"id"`
	Project   int64    `json:"project"`
	StartDate NullDate `json:"start_date" db:"start_date"`
	EndDate   NullDate `json:"end_date" db:"end_date"`
}

// ProgramCreateRequest is the data a user should submit to create a program
type ProgramCreateRequest struct {
	Project   int64    `json:"project"`
	StartDate NullDate `json:"start_date" db:"start_date" schema:"start_date"`
	EndDate   NullDate `json:"end_date" db:"end_date" schema:"end_date"`
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
