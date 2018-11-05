package main

import "github.com/paulmach/orb/encoding/wkt"

// DatapointRepository is the set of methods available for interacting with Datapoint records
type DatapointRepository interface {
	ListDatapoints() ([]*Datapoint, error)
	CreateDatapoint(dp Datapoint) (Datapoint, error)
	GetDatapoint(datapointID int) (Datapoint, error)
}

// Datapoint database methods

// CreateDatapoint creates a datapoint record.
// It may be called while handling create requests for boreholes or instruments
func (db *Datastore) CreateDatapoint(dp Datapoint) (Datapoint, error) {
	query := `INSERT INTO datapoint (location) VALUES ($1) RETURNING id`
	created := Datapoint{}
	err := db.Get(&created, query, wkt.MarshalString(dp.Location))
	if err != nil {
		return Datapoint{}, err
	}

	return created, nil
}
