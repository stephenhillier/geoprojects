package repository

import (
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/stephenhillier/geoprojects/earthworks"
	"github.com/stephenhillier/geoprojects/earthworks/db"
)

// DatapointRepository is the set of methods available for interacting with Datapoint records
type DatapointRepository interface {
	// ListDatapoints() ([]*Datapoint, error)
	CreateDatapoint(dp earthworks.Datapoint) (earthworks.Datapoint, error)
	// GetDatapoint(datapointID int) (earthworks.Datapoint, error)
}

// NewDatapointRepo returns a PostgresDatapointRepo with a database connection
func NewDatapointRepo(database *db.Datastore) *PostgresDatapointRepo {
	return &PostgresDatapointRepo{
		database,
	}
}

// PostgresDatapointRepo has a database connection and methods to interact with datapoints in
// the database.
type PostgresDatapointRepo struct {
	*db.Datastore
}

// CreateDatapoint creates a datapoint record.
// It may be called while handling create requests for boreholes or instruments
func (db *PostgresDatapointRepo) CreateDatapoint(dp earthworks.Datapoint) (earthworks.Datapoint, error) {
	query := `INSERT INTO datapoint (location) VALUES ($1) RETURNING id`
	created := earthworks.Datapoint{}
	err := db.Get(&created, query, wkt.MarshalString(dp.Location))
	if err != nil {
		return earthworks.Datapoint{}, err
	}

	return created, nil
}
