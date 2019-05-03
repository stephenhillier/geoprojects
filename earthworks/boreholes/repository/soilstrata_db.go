package repository

import (
	"github.com/stephenhillier/geoprojects/earthworks"
)

// ListStrataByBorehole retrieves a list of soil strata records associated with a given borehole
func (repo *PostgresRepo) ListStrataByBorehole(boreholeID int64) ([]*earthworks.Strata, error) {
	query := `SELECT id, borehole, start_depth, end_depth, description, soils, moisture, consistency FROM strata WHERE borehole=$1 ORDER BY start_depth`

	var err error
	strata := []*earthworks.Strata{}

	err = repo.conn.Select(&strata, query, boreholeID)

	if err != nil {
		return []*earthworks.Strata{}, err
	}
	return strata, nil
}

// CreateStrata inserts a strata record into the datastore
func (repo *PostgresRepo) CreateStrata(strata earthworks.Strata) (earthworks.Strata, error) {
	query := `
		INSERT INTO strata (borehole, start_depth, end_depth, description, soils, moisture, consistency)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, borehole, start_depth, end_depth, description, soils, moisture, consistency
	`
	created := earthworks.Strata{}
	err := repo.conn.Get(
		&created,
		query,
		strata.Borehole,
		strata.Start,
		strata.End,
		strata.Description,
		strata.Soils,
		strata.Moisture,
		strata.Consistency,
	)
	if err != nil {
		return earthworks.Strata{}, err
	}
	return created, nil
}

// CountStrataForBorehole returns a count of all strata (soil layers) in a given borehole
func (repo *PostgresRepo) CountStrataForBorehole(boreholeID int64) (int, error) {
	var count int
	query := `SELECT count(*) FROM strata WHERE borehole=$1`
	err := repo.conn.Get(&count, query, boreholeID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// RetrieveStrata gets a single strata record from the database
func (repo *PostgresRepo) RetrieveStrata(strataID int) (earthworks.Strata, error) {
	strata := earthworks.Strata{}
	query := `SELECT id, borehole, start_depth, end_depth, description, soils, moisture, consistency FROM strata WHERE id = $1`
	err := repo.conn.Get(&strata, query, strataID)
	if err != nil {
		return earthworks.Strata{}, err
	}

	return strata, nil
}

// UpdateStrata updates a Strata record in the database
func (repo *PostgresRepo) UpdateStrata(strata earthworks.Strata) (earthworks.Strata, error) {
	query := `
	UPDATE strata
	SET
		borehole = $1,
		start_depth = $2,
		end_depth = $3,
		description = $4,
		soils = $5,
		moisture = $6,
		consistency = $7
	WHERE strata.id = $8
	RETURNING id, borehole, start_depth, end_depth, description, soils, moisture, consistency
`
	created := earthworks.Strata{}
	err := repo.conn.Get(
		&created,
		query,
		strata.Borehole,
		strata.Start,
		strata.End,
		strata.Description,
		strata.Soils,
		strata.Moisture,
		strata.Consistency,
		strata.ID,
	)
	if err != nil {
		return earthworks.Strata{}, err
	}
	return created, nil
}

// DeleteStrata deletes a strata record with a given ID
func (repo *PostgresRepo) DeleteStrata(strataID int64) error {
	query := `DELETE from strata WHERE id = $1`

	_, err := repo.conn.Exec(query, strataID)
	return err
}
