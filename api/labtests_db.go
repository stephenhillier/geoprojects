package main

import (
	"errors"
	"log"
)

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
			lab_test.checked_by,
			lab_test.checked_date,
			soil_sample.name AS sample_name,
			borehole.id AS borehole,
			borehole.name AS borehole_name
		FROM lab_test
		LEFT JOIN soil_sample ON (lab_test.sample = soil_sample.id)
		LEFT JOIN borehole ON (soil_sample.borehole = borehole.id)
		WHERE borehole.project=$1
	`

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

	tx, err := db.Beginx()
	if err != nil {
		return LabTestResponse{}, err
	}

	err = tx.Get(
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
		tx.Rollback()
		return LabTestResponse{}, err
	}

	switch labTest.Type {
	case "moisture_content":
		testQuery := `
				INSERT INTO moisture_test (id) VALUES ($1)
			`
		_, err = tx.Exec(testQuery, created.ID)
		if err != nil {
			tx.Rollback()
			return LabTestResponse{}, err
		}
	default:
		tx.Rollback()
		return LabTestResponse{}, errors.New("Test type not implemented")
	}
	tx.Commit()

	return created, nil
}

// CreateMoistureTest creates a moisture_test record and updates the corresponding
// lab_test record (if necessary)
func (db *Datastore) CreateMoistureTest(labTest MoistureTestRequest, testID int) (MoistureTestResponse, error) {
	query := `
		INSERT INTO moisture_test (id, tare_mass, sample_plus_tare, dry_plus_tare)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var createdID int
	err := db.Get(
		&createdID,
		query,
		testID,
		labTest.TareMass,
		labTest.SamplePlusTare,
		labTest.DryPlusTare,
	)

	if err != nil {
		return MoistureTestResponse{}, err
	}

	created, err := db.RetrieveMoistureTest(createdID)
	if err != nil {
		return MoistureTestResponse{}, err
	}

	return created, nil
}

// UpdateMoistureTest updates a moisture_test record
func (db *Datastore) UpdateMoistureTest(labTest MoistureTestRequest, testID int) (MoistureTestResponse, error) {

	query := `
		WITH up_moisture_test AS (
			UPDATE moisture_test
			SET
				tare_mass = $1,
				sample_plus_tare = $2,
				dry_plus_tare = $3
			WHERE id = $9
			RETURNING id, tare_mass, sample_plus_tare, dry_plus_tare
		),
		up_lab_test AS (
			UPDATE lab_test
			SET
				start_date = $4,
				end_date = $5,
				performed_by = $6,
				checked_by = $7,
				checked_date = $8
			WHERE id = $9
			RETURNING id, name, start_date, end_date, performed_by, type, sample, checked_by, checked_date
		)
		SELECT
			up_lab_test.id,
			up_lab_test.name,
			up_lab_test.type,
			up_lab_test.start_date,
			up_lab_test.end_date,
			up_lab_test.performed_by,
			up_lab_test.sample,
			up_lab_test.checked_by,
			up_lab_test.checked_date,
			soil_sample.name AS sample_name,
			up_moisture_test.tare_mass,
			up_moisture_test.sample_plus_tare,
			up_moisture_test.dry_plus_tare,
			borehole.id AS borehole,
			borehole.name AS borehole_name
		FROM up_lab_test
		LEFT JOIN soil_sample ON (up_lab_test.sample = soil_sample.id)
		LEFT JOIN up_moisture_test ON (up_lab_test.id = up_moisture_test.id)
		LEFT JOIN borehole ON (soil_sample.borehole = borehole.id)
		WHERE up_lab_test.id = $9
	`

	updated := MoistureTestResponse{}

	err := db.Get(
		&updated,
		query,
		labTest.TareMass,
		labTest.SamplePlusTare,
		labTest.DryPlusTare,
		labTest.StartDate,
		labTest.EndDate,
		labTest.PerformedBy,
		labTest.CheckedBy,
		labTest.EndDate,
		testID,
	)
	if err != nil {
		return MoistureTestResponse{}, err
	}

	return updated, nil
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

// RetrieveMoistureTest gets a moisture test from the database, including relevant info
// from the lab_test table.
func (db *Datastore) RetrieveMoistureTest(moistureTestID int) (MoistureTestResponse, error) {
	moistureTest := MoistureTestResponse{}
	query := `
		SELECT
			lab_test.id,
			lab_test.name,
			lab_test.type,
			lab_test.start_date,
			lab_test.end_date,
			lab_test.performed_by,
			lab_test.sample,
			lab_test.checked_by,
			lab_test.checked_date,
			soil_sample.name AS sample_name,
			moisture_test.tare_mass,
			moisture_test.sample_plus_tare,
			moisture_test.dry_plus_tare,
			borehole.id AS borehole,
			borehole.name AS borehole_name
		FROM lab_test
		LEFT JOIN soil_sample ON (lab_test.sample = soil_sample.id)
		LEFT JOIN moisture_test ON (lab_test.id = moisture_test.id)
		LEFT JOIN borehole ON (soil_sample.borehole = borehole.id)
		WHERE lab_test.id = $1`
	err := db.Get(&moistureTest, query, moistureTestID)
	if err != nil {
		return MoistureTestResponse{}, err
	}

	return moistureTest, nil
}

// UpdateLabTest updates a test record in the database
func (db *Datastore) UpdateLabTest(labTest LabTest) (LabTestResponse, error) {
	query := `
	UPDATE lab_test
	SET
		name = $1,
		start_date = $2,
		end_date = $3,
		performed_by = $4,
		checked_date = $5,
		checked_by = $6

	WHERE id = $7
	RETURNING id, name, type, start_date, end_date, performed_by, checked_date, checked_by, sample
`
	created := LabTestResponse{}
	err := db.Get(
		&created,
		query,
		labTest.Name,
		labTest.StartDate,
		labTest.EndDate,
		labTest.PerformedBy,
		labTest.CheckedDate,
		labTest.CheckedBy,
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
