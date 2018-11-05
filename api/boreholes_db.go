package main

import "github.com/paulmach/orb"

// BoreholeRepository is the set of methods available for interacting with Borehole records
type BoreholeRepository interface {
	ListBoreholes(projectID int) ([]*BoreholeResponse, error)
	CreateBorehole(bh BoreholeCreateRequest) (Borehole, error)
	GetBorehole(boreholeID int) (Borehole, error)
}

// Borehole database methods

// ListBoreholes returns all boreholes, or, with optional projectID,
// all boreholes for a given project.
func (db *Datastore) ListBoreholes(projectID int) ([]*BoreholeResponse, error) {
	query := `SELECT id, project, program, datapoint, start_date, end_date, field_eng FROM borehole`

	queryByProject := `
		SELECT id, project, program, datapoint, name, start_date, end_date, field_eng
		FROM borehole WHERE project=$1
	`

	var err error
	boreholes := []*BoreholeResponse{}
	if projectID == 0 {
		err = db.Select(&boreholes, query)
	} else {
		err = db.Select(&boreholes, queryByProject, projectID)
	}

	if err != nil {
		return []*BoreholeResponse{}, err
	}
	return boreholes, nil
}

// CreateBorehole creates a borehole record, as well as a Datapoint record if an existing
// datapoint wasn't supplied.
// Either a datapoint or a location should be supplied.
func (db *Datastore) CreateBorehole(bh BoreholeCreateRequest) (Borehole, error) {

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

	query := `
		INSERT INTO borehole (datapoint, program, project, name, start_date, end_date, field_eng)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, project, program, name, datapoint, start_date, end_date, field_eng
	`

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
func (db *Datastore) GetBorehole(boreholeID int) (Borehole, error) {
	p := Borehole{}
	query := `SELECT id, project, program, datapoint, name, start_date, end_date, field_eng FROM borehole WHERE id=$1`
	err := db.Get(&p, query, boreholeID)
	if err != nil {
		return Borehole{}, err
	}
	return p, nil
}
