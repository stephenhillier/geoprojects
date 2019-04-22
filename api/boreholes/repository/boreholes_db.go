package repository

import (
	"log"

	"github.com/paulmach/orb"
	boreholev1 "github.com/stephenhillier/geoprojects/api/boreholes/model"
	database "github.com/stephenhillier/geoprojects/api/db"
	projectsv1 "github.com/stephenhillier/geoprojects/api/projects/model"
	projectsRepo "github.com/stephenhillier/geoprojects/api/projects/repository"
)

// BoreholeRepository is the set of methods available for interacting with Borehole records
type BoreholeRepository interface {
	ListBoreholes(int, int, int) ([]*boreholev1.BoreholeResponse, int64, error)
	CreateBorehole(bh boreholev1.BoreholeCreateRequest, project int64) (boreholev1.Borehole, error)
	GetBorehole(boreholeID int) (boreholev1.BoreholeResponse, error)
	DeleteBorehole(boreholeID int64) error

	ListStrataByBorehole(boreholeID int64) ([]*boreholev1.Strata, error)
	CreateStrata(strata boreholev1.Strata) (boreholev1.Strata, error)
	CountStrataForBorehole(boreholeID int64) (int, error)
	RetrieveStrata(strataID int) (boreholev1.Strata, error)
	UpdateStrata(strata boreholev1.Strata) (boreholev1.Strata, error)
	DeleteStrata(strataID int64) error
}

// NewBoreholeRepo returns a PostgresRepo with a database connection
func NewBoreholeRepo(database *database.Datastore) *PostgresRepo {
	return &PostgresRepo{
		database,
	}
}

// PostgresRepo has a database connection and methods to interact with boreholes in
// a Postgres database.
type PostgresRepo struct {
	conn *database.Datastore
}

// Borehole database methods

// ListBoreholes returns all boreholes, or, with optional projectID,
// all boreholes for a given project.
func (db *PostgresRepo) ListBoreholes(projectID int, limit int, offset int) ([]*boreholev1.BoreholeResponse, int64, error) {
	countQuery := `SELECT count(id) FROM borehole`
	countByProjectQuery := `SELECT count(id) FROM borehole WHERE project=$1`

	query := `
		SELECT borehole.id, borehole.project, borehole.program, borehole.datapoint, borehole.name, borehole.start_date, borehole.end_date, borehole.field_eng,
			ST_AsBinary(datapoint.location) AS location
		FROM borehole
		LEFT JOIN datapoint ON (datapoint.id = borehole.datapoint)
		LIMIT $1 OFFSET $2
	`
	queryByProject := `
		SELECT borehole.id, borehole.project, borehole.program, borehole.datapoint, borehole.name, borehole.start_date, borehole.end_date, borehole.field_eng,
			ST_AsBinary(datapoint.location) AS location
		FROM borehole
		LEFT JOIN datapoint ON (datapoint.id = borehole.datapoint)
		WHERE project=$1
		LIMIT $2 OFFSET $3
	`

	var err error
	var count int64
	boreholes := []*boreholev1.BoreholeResponse{}

	// Get counts from database
	// queries are split up this way to handle errors one at a time (counts then select queries)
	if projectID == 0 {
		err = db.conn.Get(&count, countQuery)
	} else {
		err = db.conn.Get(&count, countByProjectQuery, projectID)
	}
	if err != nil {
		// count failed
		return []*boreholev1.BoreholeResponse{}, 0, err
	}

	// select boreholes from DB
	if projectID == 0 {
		err = db.conn.Select(&boreholes, query, limit, offset)
	} else {
		err = db.conn.Select(&boreholes, queryByProject, projectID, limit, offset)
	}

	if err != nil {
		// borehole query failed

		log.Println(err)
		return []*boreholev1.BoreholeResponse{}, 0, err
	}
	return boreholes, count, nil
}

// CreateBorehole creates a borehole record, as well as a Datapoint record if an existing
// datapoint wasn't supplied.
// Either a datapoint or a location should be supplied.
func (db *PostgresRepo) CreateBorehole(bh boreholev1.BoreholeCreateRequest, project int64) (boreholev1.Borehole, error) {

	// If a datapoint wasn't supplied, create one.
	// If a location also wasn't supplied, it will be created at the default location (0, 0?)
	if !bh.Datapoint.Valid {
		projectsRepo := projectsRepo.NewDatapointRepo(db.conn)
		newDP := projectsv1.Datapoint{Location: orb.Point{bh.Location[0], bh.Location[1]}}
		createdDP, err := projectsRepo.CreateDatapoint(newDP)
		if err != nil {
			return boreholev1.Borehole{}, err
		}
		bh.Datapoint = createdDP.ID
	}

	query := `
		INSERT INTO borehole (datapoint, program, project, name, start_date, end_date, field_eng, drilling_method, type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, project, program, name, datapoint, start_date, end_date, field_eng
	`

	created := boreholev1.Borehole{}
	err := db.conn.Get(
		&created,
		query,
		bh.Datapoint,
		bh.Program,
		project,
		bh.Name,
		bh.StartDate,
		bh.EndDate,
		bh.FieldEng,
		bh.DrillingMethod,
		bh.Type,
	)
	if err != nil {
		return boreholev1.Borehole{}, err
	}

	return created, nil
}

// GetBorehole retrieves a single borehole record.
func (db *PostgresRepo) GetBorehole(boreholeID int) (boreholev1.BoreholeResponse, error) {
	p := boreholev1.BoreholeResponse{}
	query := `
		SELECT
			borehole.id,
			borehole.project,
			borehole.program,
			borehole.datapoint,
			borehole.name,
			borehole.start_date,
			borehole.end_date,
			borehole.field_eng,
			borehole.drilling_method,
			bh_drilling_method.description AS drilling_method_description,
			ST_AsBinary(datapoint.location) AS location
		FROM borehole
		JOIN drilling_method AS bh_drilling_method ON (borehole.drilling_method = bh_drilling_method.code)
		LEFT JOIN datapoint ON (datapoint.id = borehole.datapoint)
		WHERE borehole.id=$1	
		`
	err := db.conn.Get(&p, query, boreholeID)
	if err != nil {
		log.Println(err)
		return boreholev1.BoreholeResponse{}, err
	}
	return p, nil
}

// DeleteBorehole deletes a borehole by its ID
func (db *PostgresRepo) DeleteBorehole(boreholeID int64) error {
	query := `DELETE from borehole WHERE id = $1`

	_, err := db.conn.Exec(query, boreholeID)
	return err
}
