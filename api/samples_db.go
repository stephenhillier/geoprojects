package main

// ListSamplesByBorehole retrieves a list of soil samples records associated with a given borehole
func (db *Datastore) ListSamplesByBorehole(boreholeID int64) ([]*Sample, error) {
	query := `SELECT id, borehole, name, start_depth, end_depth, description, uscs FROM soil_sample WHERE borehole=$1 ORDER BY start_depth`

	var err error
	sample := []*Sample{}

	err = db.Select(&sample, query, boreholeID)

	if err != nil {
		return []*Sample{}, err
	}
	return sample, nil
}

// CreateSample inserts a sample record into the datastore
func (db *Datastore) CreateSample(sample Sample) (Sample, error) {
	query := `
		INSERT INTO soil_sample (borehole, name, start_depth, end_depth, description, uscs)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, borehole, start_depth, end_depth, description, uscs
	`
	created := Sample{}
	err := db.Get(
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
		return Sample{}, err
	}
	return created, nil
}

// CountSampleForBorehole returns a count of all sample (soil layers) in a given borehole
func (db *Datastore) CountSampleForBorehole(boreholeID int64) (int, error) {
	var count int
	query := `SELECT count(*) FROM soil_sample WHERE borehole=$1`
	err := db.Get(&count, query, boreholeID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// RetrieveSample gets a single sample record from the database
func (db *Datastore) RetrieveSample(sampleID int) (Sample, error) {
	sample := Sample{}
	query := `SELECT id, borehole, name, start_depth, end_depth, description, uscs FROM soil_sample WHERE id = $1`
	err := db.Get(&sample, query, sampleID)
	if err != nil {
		return Sample{}, err
	}

	return sample, nil
}

// UpdateSample updates a Sample record in the database
func (db *Datastore) UpdateSample(sample Sample) (Sample, error) {
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
	created := Sample{}
	err := db.Get(
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
		return Sample{}, err
	}
	return created, nil
}

// DeleteSample deletes a sample record with a given ID
func (db *Datastore) DeleteSample(sampleID int64) error {
	query := `DELETE from soil_sample WHERE id = $1`

	_, err := db.Exec(query, sampleID)
	return err
}

// CountTestForSample returns a count of all sample (soil layers) in a given sample
func (db *Datastore) CountTestForSample(sampleID int64) (int, error) {
	var count int
	query := `SELECT count(*) FROM soil_sample WHERE sample=$1`
	err := db.Get(&count, query, sampleID)
	if err != nil {
		return 0, err
	}
	return count, nil
}
