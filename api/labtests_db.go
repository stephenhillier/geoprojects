package main

import "log"

// ListLabTestsByProject retrieves a list of soil lab tests records associated with a given borehole
func (db *Datastore) ListLabTestsByProject(projectID int, boreholeID int) ([]*LabTestResponse, error) {
	query := `
		SELECT
			lab_test.id,
			lab_test.sample,
			lab_test.name,
			lab_test.start_date,
			lab_test.end_date,
			lab_test.type,
			lab_test.performed_by,
			soil_sample.name AS sample_name,
			borehole.id AS borehole,
			borehole.name AS borehole_name
		FROM lab_test
		LEFT JOIN soil_sample ON (lab_test.sample = soil_sample.id)
		LEFT JOIN borehole ON (soil_sample.borehole = borehole.id)
		WHERE borehole.project=$1
	`

	if boreholeID > 0 {
	}

	var err error
	test := []*LabTestResponse{}

	if boreholeID > 0 {
		query = query + ` AND borehole.id = $2`
		err = db.Select(&test, query, projectID, boreholeID)
	} else {
		err = db.Select(&test, query, projectID)
	}

	if err != nil {
		log.Println(err)
		return []*LabTestResponse{}, err
	}
	return test, nil
}

// CreateLabTest creates a lab test record in the datastore
func (db *Datastore) CreateLabTest(labTest LabTest) (LabTestResponse, error) {
	query := `
		INSERT INTO lab_test (name, type, start_date, end_date, performed_by, sample)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, type, start_date, end_date, performed_by, sample
	`
	created := LabTestResponse{}
	err := db.Get(
		&created,
		query,
		labTest.Name,
		labTest.Type,
		labTest.StartDate,
		labTest.EndDate,
		labTest.PerformedBy,
		labTest.Sample,
	)
	if err != nil {
		return LabTestResponse{}, err
	}
	return created, nil
}

// RetrieveLabTest gets a single lab test record from the database
func (db *Datastore) RetrieveLabTest(labTestID int) (LabTestResponse, error) {
	labTest := LabTestResponse{}
	query := `
		SELECT
			lab_test.id,
			lab_test.name,
			lab_test.type,
			lab_test.start_date,
			lab_test.end_date,
			lab_test.performed_by,
			lab_test.sample,
			soil_sample.name AS sample_name
		FROM lab_test
		LEFT JOIN soil_sample ON (lab_test.sample = soil_sample.id)
		WHERE lab_test.id = $1`
	err := db.Get(&labTest, query, labTestID)
	if err != nil {
		return LabTestResponse{}, err
	}

	return labTest, nil
}

// UpdateLabTest updates a test record in the database
func (db *Datastore) UpdateLabTest(labTest LabTest) (LabTestResponse, error) {
	query := `
	UPDATE lab_test
	SET
		name = $1,
		type = $2,
		start_date = $3,
		end_date = $4,
		performed_by = $5,
		sample = $6
	WHERE id = $7
	RETURNING id, name, type, start_date, end_date, performed_by, sample
`
	created := LabTestResponse{}
	err := db.Get(
		&created,
		query,
		labTest.Name,
		labTest.Type,
		labTest.StartDate,
		labTest.EndDate,
		labTest.PerformedBy,
		labTest.Sample,
		labTest.ID,
	)
	if err != nil {
		return LabTestResponse{}, err
	}
	return created, nil
}

// DeleteLabTest deletes a lab test record with a given ID
func (db *Datastore) DeleteLabTest(testID int) error {
	query := `DELETE from lab_test WHERE id = $1`

	_, err := db.Exec(query, testID)
	return err
}
