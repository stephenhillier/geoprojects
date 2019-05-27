package earthworks

import (
	"github.com/stephenhillier/geoprojects/earthworks/db"
)

// LabTest represents one lab test ordered on a given sample. A lab test should have
// at least one record of a specific test referencing it (e.g. a moisture content test record)
type LabTest struct {
	ID          int           `json:"id"`
	Name        db.NullString `json:"name"`
	Type        string        `json:"test_type" db:"type"`
	Sample      int           `json:"sample" db:"sample"`
	StartDate   db.NullDate   `json:"start_date" db:"start_date"`
	EndDate     db.NullDate   `json:"end_date" db:"end_date"`
	PerformedBy db.NullString `json:"performed_by" db:"performed_by"`
	CheckedDate db.NullDate   `json:"checked_date" db:"checked_date"`
	CheckedBy   db.NullString `json:"checked_by" db:"checked_by"`
}

// LabTestResponse is the data returned by the API containing info about a lab test.
type LabTestResponse struct {
	ID           int           `json:"id"`
	Name         db.NullString `json:"name"`
	Type         string        `json:"test_type" db:"type"`
	Sample       int           `json:"sample" db:"sample"`
	Borehole     int           `json:"borehole" db:"borehole"`
	BoreholeName string        `json:"borehole_name" db:"borehole_name"`
	StartDate    db.NullDate   `json:"start_date" db:"start_date"`
	EndDate      db.NullDate   `json:"end_date" db:"end_date"`
	PerformedBy  db.NullString `json:"performed_by" db:"performed_by"`
	CheckedDate  db.NullDate   `json:"checked_date" db:"checked_date"`
	CheckedBy    db.NullString `json:"checked_by" db:"checked_by"`
	SampleName   string        `json:"sample_name" db:"sample_name"`
}

// MoistureTestRequest is the data needed to create a moisture test record
type MoistureTestRequest struct {
	StartDate      db.NullDate   `json:"start_date" db:"start_date"`
	EndDate        db.NullDate   `json:"end_date" db:"end_date"`
	PerformedBy    db.NullString `json:"performed_by" db:"performed_by"`
	CheckedDate    db.NullDate   `json:"checked_date" db:"checked_date"`
	CheckedBy      db.NullString `json:"checked_by" db:"checked_by"`
	TareMass       *float64      `json:"tare_mass,string" db:"tare_mass"`
	SamplePlusTare *float64      `json:"sample_plus_tare,string" db:"sample_plus_tare"`
	DryPlusTare    *float64      `json:"dry_plus_tare,string" db:"dry_plus_tare"`
}

// MoistureTestResponse contains all the data relating to a moisture content test
type MoistureTestResponse struct {
	// ID is the ID of the lab_test and moisture_test (1 : 0..1)
	ID             int           `json:"id"`
	Name           db.NullString `json:"name"`
	Type           string        `json:"test_type" db:"type"`
	Sample         int           `json:"sample" db:"sample"`
	Borehole       int           `json:"borehole" db:"borehole"`
	BoreholeName   string        `json:"borehole_name" db:"borehole_name"`
	StartDate      db.NullDate   `json:"start_date" db:"start_date"`
	EndDate        db.NullDate   `json:"end_date" db:"end_date"`
	PerformedBy    db.NullString `json:"performed_by" db:"performed_by"`
	CheckedDate    db.NullDate   `json:"checked_date" db:"checked_date"`
	CheckedBy      db.NullString `json:"checked_by" db:"checked_by"`
	SampleName     string        `json:"sample_name" db:"sample_name"`
	TareMass       *float64      `json:"tare_mass" db:"tare_mass"`
	SamplePlusTare *float64      `json:"sample_plus_tare" db:"sample_plus_tare"`
	DryPlusTare    *float64      `json:"dry_plus_tare" db:"dry_plus_tare"`
}

// GSATestRequest is the set of data that the client provides to
// start or update a grain size analysis test record
type GSATestRequest struct {
	StartDate      db.NullDate      `json:"start_date" db:"start_date"`
	EndDate        db.NullDate      `json:"end_date" db:"end_date"`
	PerformedBy    db.NullString    `json:"performed_by" db:"performed_by"`
	CheckedDate    db.NullDate      `json:"checked_date" db:"checked_date"`
	CheckedBy      db.NullString    `json:"checked_by" db:"checked_by"`
	TareMass       *float64         `json:"tare_mass,string" db:"tare_mass"`
	SamplePlusTare *float64         `json:"sample_plus_tare,string" db:"sample_plus_tare"`
	WashedPlusTare *float64         `json:"washed_plus_tare,string" db:"washed_plus_tare"`
	DryPlusTare    *float64         `json:"dry_plus_tare,string" db:"dry_plus_tare"`
	Sieves         []GSADataRequest `json:"sieves"`
}

// GSATestResponse is the set of data that the client provides to
// start or update a grain size analysis test record
type GSATestResponse struct {
	ID             int                `json:"id"`
	Name           db.NullString      `json:"name"`
	Type           string             `json:"test_type" db:"type"`
	Sample         int                `json:"sample" db:"sample"`
	SampleName     string             `json:"sample_name" db:"sample_name"`
	Borehole       int                `json:"borehole" db:"borehole"`
	BoreholeName   string             `json:"borehole_name" db:"borehole_name"`
	StartDate      db.NullDate        `json:"start_date" db:"start_date"`
	EndDate        db.NullDate        `json:"end_date" db:"end_date"`
	PerformedBy    db.NullString      `json:"performed_by" db:"performed_by"`
	CheckedDate    db.NullDate        `json:"checked_date" db:"checked_date"`
	CheckedBy      db.NullString      `json:"checked_by" db:"checked_by"`
	TareMass       *float64           `json:"tare_mass" db:"tare_mass"`
	SamplePlusTare *float64           `json:"sample_plus_tare" db:"sample_plus_tare"`
	WashedPlusTare *float64           `json:"washed_plus_tare" db:"washed_plus_tare"`
	DryPlusTare    *float64           `json:"dry_plus_tare" db:"dry_plus_tare"`
	Sieves         []*GSADataResponse `json:"sieves"`
}

// GSADataRequest is the data required to add or update the test result
// from a single sieve/pan in a grain size analysis test
type GSADataRequest struct {
	// Test    int     `json:"gsa_test"`
	Pan      bool    `json:"pan"`
	Size     float64 `json:"size,string"`
	Retained float64 `json:"mass_retained,string" db:"mass_retained"`
}

// GSADataResponse is the data returned by the API when a client
// requests grain size analysis data for a single sieve.
// This is also used for a GSATestResponse as a collection of all
// the sieves for a single test.
type GSADataResponse struct {
	ID       int64   `json:"id"`
	Test     int     `json:"gsa_test"`
	Pan      bool    `json:"pan"`
	Size     float64 `json:"size"`
	Retained float64 `json:"mass_retained" db:"mass_retained"`
}

// LabTestCtx is a key for a lab test stored in context
var LabTestCtx struct{}

// LabRepository is the set of methods available for interacting with lab tests
type LabRepository interface {
	DeleteLabTest(testID int) error
	UpdateLabTest(LabTest) (LabTestResponse, error)
	RetrieveLabTest(labTestID int) (LabTestResponse, error)
	CreateLabTest(labTest LabTest) (LabTestResponse, error)
	ListLabTestsByProject(projectID int, boreholeID int) ([]*LabTestResponse, error)
	SieveTestRepository
	MoistureTestRepository
}

// SieveTestRepository has methods for working with sieve/GSA tests
type SieveTestRepository interface {
	AddSieve(test GSADataRequest, testID int) (GSADataResponse, error)
	RetrieveSieve(testID int) (GSADataResponse, error)
	RetrieveSieves(testID int) ([]*GSADataResponse, error)
	UpdateSieve(test GSADataRequest, testID int, sieveID int) (GSADataResponse, error)
	DeleteSieve(sieveID int) error
	RetrieveSieveTest(testID int) (GSATestResponse, error)
	UpdateGSATest(labTest GSATestRequest, testID int) (GSATestResponse, error)
}

// MoistureTestRepository has methods for working with moisture content tests
type MoistureTestRepository interface {
	RetrieveMoistureTest(moistureTestID int) (MoistureTestResponse, error)
	UpdateMoistureTest(data MoistureTestRequest, testID int) (MoistureTestResponse, error)
	CreateMoistureTest(labTest MoistureTestRequest, testID int) (MoistureTestResponse, error)
}
