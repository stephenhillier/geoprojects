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
func (db *Datastore) ListBoreholes(projectID int, limit int, offset int) ([]*BoreholeResponse, int64, error) {
	countQuery := `SELECT count(id) FROM borehole`
	countByProjectQuery := `SELECT count(id) FROM borehole WHERE project=$1`

	query := `SELECT id, project, program, datapoint, start_date, end_date, field_eng FROM borehole LIMIT $1 OFFSET $2`

	queryByProject := `
		SELECT id, project, program, datapoint, name, start_date, end_date, field_eng
		FROM borehole WHERE project=$1
		LIMIT $2 OFFSET $3
	`

	var err error
	var count int64
	boreholes := []*BoreholeResponse{}

	// Get counts from database
	// queries are split up this way to handle errors one at a time (counts then select queries)
	if projectID == 0 {
		err = db.Get(&count, countQuery)
	} else {
		err = db.Get(&count, countByProjectQuery, projectID)
	}
	if err != nil {
		// count failed
		return []*BoreholeResponse{}, 0, err
	}

	// select boreholes from DB
	if projectID == 0 {
		err = db.Select(&boreholes, query, limit, offset)
	} else {
		err = db.Select(&boreholes, queryByProject, projectID, limit, offset)
	}

	if err != nil {
		// borehole query failed
		return []*BoreholeResponse{}, 0, err
	}
	return boreholes, count, nil
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
