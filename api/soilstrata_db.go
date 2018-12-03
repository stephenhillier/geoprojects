package main

// ListStrataByBorehole retrieves a list of soil strata records associated with a given borehole
func (db *Datastore) ListStrataByBorehole(boreholeID int64) ([]*Strata, error) {
	query := `SELECT id, borehole, start_depth, end_depth, description, soils, moisture, consistency FROM strata WHERE borehole=$1 ORDER BY start_depth`

	var err error
	strata := []*Strata{}

	err = db.Select(&strata, query, boreholeID)

	if err != nil {
		return []*Strata{}, err
	}
	return strata, nil
}

// CreateStrata inserts a strata record into the datastore
func (db *Datastore) CreateStrata(strata Strata) (Strata, error) {
	query := `
		INSERT INTO strata (borehole, start_depth, end_depth, description, soils, moisture, consistency)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, borehole, start_depth, end_depth, description, soils, moisture, consistency
	`
	created := Strata{}
	err := db.Get(
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
		return Strata{}, err
	}
	return created, nil
}

// CountStrataForBorehole returns a count of all strata (soil layers) in a given borehole
func (db *Datastore) CountStrataForBorehole(boreholeID int64) (int, error) {
	var count int
	query := `SELECT count(*) FROM strata WHERE borehole=$1`
	err := db.Get(&count, query, boreholeID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// RetrieveStrata gets a single strata record from the database
func (db *Datastore) RetrieveStrata(strataID int) (Strata, error) {
	strata := Strata{}
	query := `SELECT id, borehole, start_depth, end_depth, description, soils, moisture, consistency FROM strata WHERE id = $1`
	err := db.Get(&strata, query, strataID)
	if err != nil {
		return Strata{}, err
	}

	return strata, nil
}

// UpdateStrata updates a Strata record in the database
func (db *Datastore) UpdateStrata(strata Strata) (Strata, error) {
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
	created := Strata{}
	err := db.Get(
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
		return Strata{}, err
	}
	return created, nil
}

// DeleteStrata deletes a strata record with a given ID
func (db *Datastore) DeleteStrata(strataID int64) error {
	query := `DELETE from strata WHERE id = $1`

	_, err := db.Exec(query, strataID)
	return err
}
