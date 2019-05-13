package repository

import (
	"log"

	"github.com/stephenhillier/geoprojects/earthworks"
)

// ListSamplesByBorehole retrieves a list of soil samples records associated with a given borehole
func (repo *PostgresRepo) ListSamplesByBorehole(boreholeID int64) ([]*earthworks.Sample, error) {
	query := `SELECT
		soil_sample.id,
		soil_sample.borehole,
		soil_sample.name,
		soil_sample.start_depth,
		soil_sample.end_depth,
		soil_sample.description,
		soil_sample.uscs,
		borehole.name AS borehole_name
		FROM soil_sample
		LEFT JOIN borehole ON (soil_sample.borehole = borehole.id)
		WHERE borehole=$1 ORDER BY start_depth`

	var err error
	sample := []*earthworks.Sample{}

	err = repo.conn.Select(&sample, query, boreholeID)

	if err != nil {
		return []*earthworks.Sample{}, err
	}
	return sample, nil
}

// ListSamplesByProject lists all samples in a project
func (repo *PostgresRepo) ListSamplesByProject(projectID int) ([]*earthworks.Sample, error) {
	query := `
		SELECT
			soil_sample.id,
			soil_sample.borehole,
			soil_sample.name,
			soil_sample.start_depth,
			soil_sample.end_depth,
			soil_sample.description,
			soil_sample.uscs,
			borehole.name AS borehole_name
		FROM soil_sample
		LEFT JOIN borehole ON (soil_sample.borehole = borehole.id)
		WHERE borehole.project=$1
		ORDER BY borehole.name, start_depth`

	var err error
	samples := []*earthworks.Sample{}

	err = repo.conn.Select(&samples, query, projectID)

	if err != nil {
		log.Println(err)
		return []*earthworks.Sample{}, err
	}
	return samples, nil
}

// CreateSample inserts a sample record into the datastore
func (repo *PostgresRepo) CreateSample(sample earthworks.Sample) (earthworks.Sample, error) {
	query := `
		INSERT INTO soil_sample (borehole, name, start_depth, end_depth, description, uscs)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, borehole, start_depth, end_depth, description, uscs
	`
	created := earthworks.Sample{}
	err := repo.conn.Get(
		&created,
		query,
		sample.Borehole,
		sample.Name,
		sample.Start,
		sample.End,
		sample.Description,
		sample.USCS,
	)
	if err != nil {
		return earthworks.Sample{}, err
	}
	return created, nil
}

// CountSampleForBorehole returns a count of all sample (soil layers) in a given borehole
func (repo *PostgresRepo) CountSampleForBorehole(boreholeID int64) (int, error) {
	var count int
	query := `SELECT count(*) FROM soil_sample WHERE borehole=$1`
	err := repo.conn.Get(&count, query, boreholeID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// RetrieveSample gets a single sample record from the database
func (repo *PostgresRepo) RetrieveSample(sampleID int) (earthworks.Sample, error) {
	sample := earthworks.Sample{}
	query := `SELECT id, borehole, name, start_depth, end_depth, description, uscs FROM soil_sample WHERE id = $1`
	err := repo.conn.Get(&sample, query, sampleID)
	if err != nil {
		return earthworks.Sample{}, err
	}

	return sample, nil
}

// UpdateSample updates a Sample record in the database
func (repo *PostgresRepo) UpdateSample(sample earthworks.Sample) (earthworks.Sample, error) {
	query := `
	UPDATE soil_sample
	SET
		borehole = $1,
		start_depth = $2,
		end_depth = $3,
		description = $4,
		uscs = $5,
		name = $6
	WHERE id = $7
	RETURNING id, name, borehole, start_depth, end_depth, description, uscs
`
	created := earthworks.Sample{}
	err := repo.conn.Get(
		&created,
		query,
		sample.Borehole,
		sample.Start,
		sample.End,
		sample.Description,
		sample.USCS,
		sample.Name,
		sample.ID,
	)
	if err != nil {
		return earthworks.Sample{}, err
	}
	return created, nil
}

// DeleteSample deletes a sample record with a given ID
func (repo *PostgresRepo) DeleteSample(sampleID int64) error {
	query := `DELETE from soil_sample WHERE id = $1`

	_, err := repo.conn.Exec(query, sampleID)
	return err
}

// CountTestForSample returns a count of all sample (soil layers) in a given sample
func (repo *PostgresRepo) CountTestForSample(sampleID int64) (int, error) {
	var count int
	query := `SELECT count(*) FROM soil_sample WHERE sample=$1`
	err := repo.conn.Get(&count, query, sampleID)
	if err != nil {
		return 0, err
	}
	return count, nil
}
