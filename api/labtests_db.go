package main

// ListLabTestsByProject retrieves a list of soil lab tests records associated with a given borehole
func (db *Datastore) ListLabTestsByProject(projectID int64) ([]*LabTest, error) {
	query := `
		SELECT
			lab_test.id,
			lab_test.sample,
			lab_test.name,
			lab_test.start_date,
			lab_test.end_date,
			lab_test.type,
			borehole.name as borehole
		FROM lab_test
		LEFT JOIN soil_sample ON (lab_test.sample = soil_sample.id)
		LEFT JOIN borehole ON (soil_sample.borehole = borehole.id)
		WHERE borehole.project=$1
	`

	var err error
	test := []*LabTest{}

	err = db.Select(&test, query, projectID)

	if err != nil {
		return []*LabTest{}, err
	}
	return test, nil
}
