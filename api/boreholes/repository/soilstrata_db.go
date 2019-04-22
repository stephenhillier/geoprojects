package repository

import (
	boreholev1 "github.com/stephenhillier/geoprojects/api/boreholes/model"
)

// ListStrataByBorehole retrieves a list of soil strata records associated with a given borehole
func (db *PostgresRepo) ListStrataByBorehole(boreholeID int64) ([]*boreholev1.Strata, error) {
	query := `SELECT id, borehole, start_depth, end_depth, description, soils, moisture, consistency FROM strata WHERE borehole=$1 ORDER BY start_depth`

	var err error
	strata := []*boreholev1.Strata{}

	err = db.conn.Select(&strata, query, boreholeID)

	if err != nil {
		return []*boreholev1.Strata{}, err
	}
	return strata, nil
}

// CreateStrata inserts a strata record into the datastore
func (db *PostgresRepo) CreateStrata(strata boreholev1.Strata) (boreholev1.Strata, error) {
	query := `
		INSERT INTO strata (borehole, start_depth, end_depth, description, soils, moisture, consistency)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, borehole, start_depth, end_depth, description, soils, moisture, consistency
	`
	created := boreholev1.Strata{}
	err := db.conn.Get(
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
		return boreholev1.Strata{}, err
	}
	return created, nil
}

// CountStrataForBorehole returns a count of all strata (soil layers) in a given borehole
func (db *PostgresRepo) CountStrataForBorehole(boreholeID int64) (int, error) {
	var count int
	query := `SELECT count(*) FROM strata WHERE borehole=$1`
	err := db.conn.Get(&count, query, boreholeID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// RetrieveStrata gets a single strata record from the database
func (db *PostgresRepo) RetrieveStrata(strataID int) (boreholev1.Strata, error) {
	strata := boreholev1.Strata{}
	query := `SELECT id, borehole, start_depth, end_depth, description, soils, moisture, consistency FROM strata WHERE id = $1`
	err := db.conn.Get(&strata, query, strataID)
	if err != nil {
		return boreholev1.Strata{}, err
	}

	return strata, nil
}

// UpdateStrata updates a Strata record in the database
func (db *PostgresRepo) UpdateStrata(strata boreholev1.Strata) (boreholev1.Strata, error) {
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
	created := boreholev1.Strata{}
	err := db.conn.Get(
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
		return boreholev1.Strata{}, err
	}
	return created, nil
}

// DeleteStrata deletes a strata record with a given ID
func (db *PostgresRepo) DeleteStrata(strataID int64) error {
	query := `DELETE from strata WHERE id = $1`

	_, err := db.conn.Exec(query, strataID)
	return err
}
