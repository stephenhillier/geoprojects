package repository

import (
	"log"

	"github.com/paulmach/orb"
	"github.com/stephenhillier/geoprojects/earthworks"
	"github.com/stephenhillier/geoprojects/earthworks/db"
	projectsRepo "github.com/stephenhillier/geoprojects/earthworks/projects/repository"
)

// BoreholeRepository is the set of methods available for interacting with Borehole records
type BoreholeRepository interface {
	ListBoreholes(int, int, int) ([]*earthworks.BoreholeResponse, int64, error)
	CreateBorehole(bh earthworks.BoreholeCreateRequest, project int64) (earthworks.Borehole, error)
	GetBorehole(boreholeID int) (earthworks.BoreholeResponse, error)
	DeleteBorehole(boreholeID int64) error

	ListStrataByBorehole(boreholeID int64) ([]*earthworks.Strata, error)
	CreateStrata(strata earthworks.Strata) (earthworks.Strata, error)
	CountStrataForBorehole(boreholeID int64) (int, error)
	RetrieveStrata(strataID int) (earthworks.Strata, error)
	UpdateStrata(strata earthworks.Strata) (earthworks.Strata, error)
	DeleteStrata(strataID int64) error
}

// NewBoreholeRepo returns a PostgresRepo with a database connection
func NewBoreholeRepo(database *db.Datastore) *PostgresRepo {
	return &PostgresRepo{
		database,
	}
}

// PostgresRepo has a database connection and methods to interact with boreholes in
// a Postgres database.
type PostgresRepo struct {
	conn *db.Datastore
}

// Borehole database methods

// ListBoreholes returns all boreholes, or, with optional projectID,
// all boreholes for a given project.
func (repo *PostgresRepo) ListBoreholes(projectID int, limit int, offset int) ([]*earthworks.BoreholeResponse, int64, error) {
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
	boreholes := []*earthworks.BoreholeResponse{}

	// Get counts from database
	// queries are split up this way to handle errors one at a time (counts then select queries)
	if projectID == 0 {
		err = repo.conn.Get(&count, countQuery)
	} else {
		err = repo.conn.Get(&count, countByProjectQuery, projectID)
	}
	if err != nil {
		// count failed
		return []*earthworks.BoreholeResponse{}, 0, err
	}

	// select boreholes from DB
	if projectID == 0 {
		err = repo.conn.Select(&boreholes, query, limit, offset)
	} else {
		err = repo.conn.Select(&boreholes, queryByProject, projectID, limit, offset)
	}

	if err != nil {
		// borehole query failed

		log.Println(err)
		return []*earthworks.BoreholeResponse{}, 0, err
	}
	return boreholes, count, nil
}

// CreateBorehole creates a borehole record, as well as a Datapoint record if an existing
// datapoint wasn't supplied.
// Either a datapoint or a location should be supplied.
func (repo *PostgresRepo) CreateBorehole(bh earthworks.BoreholeCreateRequest, project int64) (earthworks.Borehole, error) {

	// If a datapoint wasn't supplied, create one.
	// If a location also wasn't supplied, it will be created at the default location (0, 0?)
	if !bh.Datapoint.Valid {
		projectsRepo := projectsRepo.NewDatapointRepo(repo.conn)
		newDP := earthworks.Datapoint{Location: orb.Point{bh.Location[0], bh.Location[1]}}
		createdDP, err := projectsRepo.CreateDatapoint(newDP)
		if err != nil {
			return earthworks.Borehole{}, err
		}
		bh.Datapoint = createdDP.ID
	}

	query := `
		INSERT INTO borehole (datapoint, program, project, name, start_date, end_date, field_eng, drilling_method, type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, project, program, name, datapoint, start_date, end_date, field_eng
	`

	created := earthworks.Borehole{}
	err := repo.conn.Get(
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
		return earthworks.Borehole{}, err
	}

	return created, nil
}

// GetBorehole retrieves a single borehole record.
func (repo *PostgresRepo) GetBorehole(boreholeID int) (earthworks.BoreholeResponse, error) {
	p := earthworks.BoreholeResponse{}
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
	err := repo.conn.Get(&p, query, boreholeID)
	if err != nil {
		log.Println(err)
		return earthworks.BoreholeResponse{}, err
	}
	return p, nil
}

// DeleteBorehole deletes a borehole by its ID
func (repo *PostgresRepo) DeleteBorehole(boreholeID int64) error {
	query := `DELETE from borehole WHERE id = $1`

	_, err := repo.conn.Exec(query, boreholeID)
	return err
}
