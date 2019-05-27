package repository

import (
	"errors"
	"log"

	"github.com/stephenhillier/geoprojects/earthworks"
	"github.com/stephenhillier/geoprojects/earthworks/db"
)

// LabRepository is the set of methods available for interacting with lab tests
type LabRepository interface {
}

// NewLabRepo returns a PostgresRepo with a database connection
func NewLabRepo(database *db.Datastore) *PostgresRepo {
	return &PostgresRepo{
		database,
	}
}

// PostgresRepo has a database connection and methods to interact with
// a Postgres database.
type PostgresRepo struct {
	conn *db.Datastore
}

// ListLabTestsByProject retrieves a list of soil lab tests records associated with a given borehole
func (db *PostgresRepo) ListLabTestsByProject(projectID int, boreholeID int) ([]*earthworks.LabTestResponse, error) {
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
	test := []*earthworks.LabTestResponse{}

	if boreholeID > 0 {
		query = query + ` AND borehole.id = $2`
		err = db.conn.Select(&test, query, projectID, boreholeID)
	} else {
		err = db.conn.Select(&test, query, projectID)
	}

	if err != nil {
		log.Println(err)
		return []*earthworks.LabTestResponse{}, err
	}
	return test, nil
}

// CreateLabTest creates a lab test record in the datastore
func (db *PostgresRepo) CreateLabTest(labTest earthworks.LabTest) (earthworks.LabTestResponse, error) {
	query := `
		INSERT INTO lab_test (name, type, start_date, end_date, performed_by, sample)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, type, start_date, end_date, performed_by, sample
	`
	created := earthworks.LabTestResponse{}

	tx, err := db.conn.Beginx()
	if err != nil {
		return earthworks.LabTestResponse{}, err
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
		return earthworks.LabTestResponse{}, err
	}

	switch labTest.Type {
	case "moisture_content":
		testQuery := `
				INSERT INTO moisture_test (id) VALUES ($1)
			`
		_, err = tx.Exec(testQuery, created.ID)
		if err != nil {
			tx.Rollback()
			return earthworks.LabTestResponse{}, err
		}
	case "grain_size_analysis":
		testQuery := `
				INSERT INTO gsa_test (id) VALUES ($1)
			`
		_, err = tx.Exec(testQuery, created.ID)
		if err != nil {
			tx.Rollback()
			return earthworks.LabTestResponse{}, err
		}
	default:
		tx.Rollback()
		return earthworks.LabTestResponse{}, errors.New("Test type not implemented")
	}
	tx.Commit()

	return created, nil
}

// CreateMoistureTest creates a moisture_test record and updates the corresponding
// lab_test record (if necessary)
func (db *PostgresRepo) CreateMoistureTest(labTest earthworks.MoistureTestRequest, testID int) (earthworks.MoistureTestResponse, error) {
	query := `
		INSERT INTO moisture_test (id, tare_mass, sample_plus_tare, dry_plus_tare)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var createdID int
	err := db.conn.Get(
		&createdID,
		query,
		testID,
		labTest.TareMass,
		labTest.SamplePlusTare,
		labTest.DryPlusTare,
	)

	if err != nil {
		return earthworks.MoistureTestResponse{}, err
	}

	created, err := db.RetrieveMoistureTest(createdID)
	if err != nil {
		return earthworks.MoistureTestResponse{}, err
	}

	return created, nil
}

// UpdateMoistureTest updates a moisture_test record
func (db *PostgresRepo) UpdateMoistureTest(labTest earthworks.MoistureTestRequest, testID int) (earthworks.MoistureTestResponse, error) {

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

	updated := earthworks.MoistureTestResponse{}

	err := db.conn.Get(
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
		return earthworks.MoistureTestResponse{}, err
	}

	return updated, nil
}

// UpdateGSATest updates a grain size analysis record
func (db *PostgresRepo) UpdateGSATest(labTest earthworks.GSATestRequest, testID int) (earthworks.GSATestResponse, error) {

	removeOldDataQuery := `
		DELETE FROM gsa_data
		WHERE test = $1
	`

	sieveDataQuery := `
		INSERT INTO gsa_data (test, pan, size, mass_retained)
		VALUES ($1, $2, $3, $4)
		RETURNING id, test, pan, size, mass_retained
	`

	updateTestSummaryQuery := `
		WITH up_gsa_test AS (
			UPDATE gsa_test
			SET
				tare_mass = $1,
				dry_plus_tare = $2,
				washed_plus_tare = $3
			WHERE id = $9
			RETURNING id, tare_mass, dry_plus_tare, washed_plus_tare
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
			up_gsa_test.tare_mass,
			up_gsa_test.washed_plus_tare,
			up_gsa_test.dry_plus_tare,
			borehole.id AS borehole,
			borehole.name AS borehole_name
		FROM up_lab_test
		LEFT JOIN soil_sample ON (up_lab_test.sample = soil_sample.id)
		LEFT JOIN up_gsa_test ON (up_lab_test.id = up_gsa_test.id)
		LEFT JOIN borehole ON (soil_sample.borehole = borehole.id)
		WHERE up_lab_test.id = $9
	`

	updated := earthworks.GSATestResponse{}

	tx, err := db.conn.Beginx()
	if err != nil {
		tx.Rollback()
		return updated, err
	}

	_, err = tx.Exec(removeOldDataQuery, testID)
	if err != nil {
		tx.Rollback()
		return updated, err
	}

	for _, s := range labTest.Sieves {
		_, err = tx.Exec(sieveDataQuery, testID, s.Pan, s.Size, s.Retained)
		if err != nil {
			tx.Rollback()
			return updated, err
		}
	}

	err = tx.Get(
		&updated,
		updateTestSummaryQuery,
		labTest.TareMass,
		labTest.DryPlusTare,
		labTest.WashedPlusTare,
		labTest.StartDate,
		labTest.EndDate,
		labTest.PerformedBy,
		labTest.CheckedBy,
		labTest.EndDate,
		testID,
	)
	if err != nil {
		tx.Rollback()
		return earthworks.GSATestResponse{}, err
	}

	tx.Commit()

	return updated, nil
}

// RetrieveLabTest gets a single lab test record from the database
func (db *PostgresRepo) RetrieveLabTest(labTestID int) (earthworks.LabTestResponse, error) {
	labTest := earthworks.LabTestResponse{}
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
	err := db.conn.Get(&labTest, query, labTestID)
	if err != nil {
		return earthworks.LabTestResponse{}, err
	}

	return labTest, nil
}

// RetrieveMoistureTest gets a moisture test from the database, including relevant info
// from the lab_test table.
func (db *PostgresRepo) RetrieveMoistureTest(moistureTestID int) (earthworks.MoistureTestResponse, error) {
	moistureTest := earthworks.MoistureTestResponse{}
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
	err := db.conn.Get(&moistureTest, query, moistureTestID)
	if err != nil {
		return earthworks.MoistureTestResponse{}, err
	}

	return moistureTest, nil
}

// UpdateLabTest updates a test record in the database
func (db *PostgresRepo) UpdateLabTest(labTest earthworks.LabTest) (earthworks.LabTestResponse, error) {
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
	created := earthworks.LabTestResponse{}
	err := db.conn.Get(
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
		return earthworks.LabTestResponse{}, err
	}
	return created, nil
}

// DeleteLabTest deletes a lab test record with a given ID
func (db *PostgresRepo) DeleteLabTest(testID int) error {
	query := `DELETE from lab_test WHERE id = $1`

	_, err := db.conn.Exec(query, testID)
	return err
}

// AddSieve adds a single sieve record, referencing a grain size test record
func (db *PostgresRepo) AddSieve(test earthworks.GSADataRequest, testID int) (earthworks.GSADataResponse, error) {
	query := `
		INSERT INTO gsa_data (test, pan, size, mass_retained)
		VALUES ($1, $2, $3, $4)
		RETURNING id, test, pan, size, mass_retained
	`

	created := earthworks.GSADataResponse{}

	err := db.conn.Get(&created, query, testID, test.Pan, test.Size, test.Retained)
	if err != nil {
		return earthworks.GSADataResponse{}, err
	}

	return created, nil
}

// RetrieveSieve fetches a single sieve record from the datastore
func (db *PostgresRepo) RetrieveSieve(testID int) (earthworks.GSADataResponse, error) {
	query := `
		SELECT (id, test, pan, size, mass_retained)
		FROM gsa_data
		WHERE id=$1
	`

	sieve := earthworks.GSADataResponse{}
	err := db.conn.Get(&sieve, query, testID)
	if err != nil {
		return earthworks.GSADataResponse{}, err
	}

	return sieve, nil
}

// RetrieveSieves fetches all sieves for a given grain size test
func (db *PostgresRepo) RetrieveSieves(testID int) ([]*earthworks.GSADataResponse, error) {
	query := `
		SELECT (id, test, pan, size, mass_retained)
		FROM gsa_data
		WHERE test=$1
	`

	sieves := []*earthworks.GSADataResponse{}

	err := db.conn.Select(&sieves, query, testID)
	if err != nil {
		return sieves, err
	}

	return sieves, nil
}

// UpdateSieve updates a single sieve record for a grain size test
func (db *PostgresRepo) UpdateSieve(test earthworks.GSADataRequest, testID int, sieveID int) (earthworks.GSADataResponse, error) {
	query := `
		UPDATE gsa_data
		SET pan, size, mass_retained
		VALUES ($1, $2, $3)
		WHERE id = $4 AND test = $5
		RETURNING id, test, pan, size, mass_retained
	`

	created := earthworks.GSADataResponse{}

	err := db.conn.Get(&created, query, test.Pan, test.Size, test.Retained, sieveID, testID)
	if err != nil {
		return earthworks.GSADataResponse{}, err
	}

	return created, nil
}

// DeleteSieve removes a single sieve record
func (db *PostgresRepo) DeleteSieve(sieveID int) error {
	query := `
		DELETE from gsa_data WHERE id=$1
	`

	_, err := db.conn.Exec(query, sieveID)
	return err
}

// RetrieveSieveTest retrieves a sieve test record
func (db *PostgresRepo) RetrieveSieveTest(testID int) (earthworks.GSATestResponse, error) {
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
		gsa_test.tare_mass,
		gsa_test.washed_plus_tare,
		gsa_test.dry_plus_tare,
		borehole.id AS borehole,
		borehole.name AS borehole_name
	FROM lab_test
	LEFT JOIN soil_sample ON (lab_test.sample = soil_sample.id)
	LEFT JOIN gsa_test ON (lab_test.id = gsa_test.id)
	LEFT JOIN borehole ON (soil_sample.borehole = borehole.id)
	WHERE lab_test.id = $1
	`

	querySieves := `
	SELECT id, test, pan, size, mass_retained
	FROM gsa_data WHERE test = $1
	`

	test := earthworks.GSATestResponse{}
	sieves := []*earthworks.GSADataResponse{}

	err := db.conn.Get(&test, query, testID)
	if err != nil {
		return test, err
	}

	err = db.conn.Select(&sieves, querySieves, testID)
	if err != nil {
		return test, err
	}

	test.Sieves = sieves

	return test, nil

}
