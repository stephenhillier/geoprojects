package repository

import (
	"github.com/paulmach/orb"
	"github.com/stephenhillier/geoprojects/earthworks"
	"github.com/stephenhillier/geoprojects/earthworks/db"
	projectsRepo "github.com/stephenhillier/geoprojects/earthworks/projects/repository"
)

// NewInstrumentationRepo returns a PostgresRepo with a database connection
// This method can be called with either a sqlx.DB or a sqlx.Tx (transaction)
func NewInstrumentationRepo(database *db.Datastore) *PostgresRepo {
	return &PostgresRepo{
		conn: database,
	}
}

// PostgresRepo has a database connection and methods to interact with instrumentation in
// the database.
type PostgresRepo struct {
	conn *db.Datastore
}

// ListInstruments lists all instruments in a given project.
// TODO: add a filter option for more search options
func (repo *PostgresRepo) ListInstruments(project int) ([]*earthworks.Instrument, error) {
	query := `
		SELECT instrument.id, instrument.project, instrument.program, instrument.datapoint, instrument.name, instrument.install_date, instrument.device_id, instrument.field_eng,
			ST_AsBinary(datapoint.location) AS location
		FROM instrument
		LEFT JOIN datapoint ON (datapoint.id = instrument.datapoint)
		WHERE project=$1
	`

	instruments := []*earthworks.Instrument{}

	err := repo.conn.Select(&instruments, query, project)
	if err != nil {
		return instruments, err
	}
	return instruments, err
}

// CreateInstrument adds an instrument to the database.
// Instruments reference a datapoint (containing a location);
// if none was supplied, one will be created.
// A datapoint will normally be supplied if this instrument was installed
// directly into a borehole, and therefore has the same location.
// Either a datapoint or a location should be supplied.
func (repo *PostgresRepo) CreateInstrument(instr earthworks.InstrumentCreateRequest, project int) (earthworks.Instrument, error) {

	// If a datapoint wasn't supplied, create one.
	// If a location also wasn't supplied, it will be created at the default location (0, 0?)
	if !instr.Datapoint.Valid {
		projectsRepo := projectsRepo.NewDatapointRepo(repo.conn)
		newDP := earthworks.Datapoint{Location: orb.Point{instr.Location[0], instr.Location[1]}}
		createdDP, err := projectsRepo.CreateDatapoint(newDP)
		if err != nil {
			return earthworks.Instrument{}, err
		}
		instr.Datapoint = createdDP.ID
	}

	query := `
		INSERT INTO instrument (project, name, device_id, datapoint, field_eng, install_date)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, project, name, device_id, datapoint, field_eng, install_date
	`

	created := earthworks.Instrument{}
	err := repo.conn.Get(
		&created,
		query,
		project,
		instr.Name,
		instr.DeviceID,
		instr.Datapoint,
		instr.FieldEng,
		instr.InstallDate,
	)
	if err != nil {
		return earthworks.Instrument{}, err
	}

	return created, nil
}

// GetInstrument retrieves details about a single instrument by its ID
func (repo *PostgresRepo) GetInstrument(id int64) (earthworks.Instrument, error) {

	query := `
		SELECT
			instrument.id,
			instrument.project,
			instrument.name,
			instrument.device_id,
			instrument.datapoint,
			instrument.field_eng,
			instrument.install_date,
			ST_AsBinary(datapoint.location) AS location
		FROM instrument
		JOIN datapoint ON instrument.datapoint = datapoint.id
		WHERE instrument.id = $1	
	`

	instr := earthworks.Instrument{}

	err := repo.conn.Get(&instr, query, id)
	return instr, err
}

// PostTimeSeriesData takes any data in an earthworks.TimeSeriesData type
// (normally containing a timestamp and a value), and stores it in the database.
// The DeviceID is normally the serial number of the device.
func (repo *PostgresRepo) PostTimeSeriesData(data earthworks.TimeSeriesData) (earthworks.TimeSeriesData, error) {
	query := `
		INSERT INTO time_series_data (device_id, time, value)
		VALUES ($1, $2, $3)
		RETURNING id, device_id, time as timestamp, value
	`

	created := earthworks.TimeSeriesData{}
	err := repo.conn.Get(&created, query, data.DeviceID, data.Timestamp, data.Value)

	return created, err
}
