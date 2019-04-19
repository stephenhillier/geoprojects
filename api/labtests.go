package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	projects_v1 "github.com/stephenhillier/geoprojects/api/projects/model"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// LabTest represents one lab test ordered on a given sample. A lab test should have
// at least one record of a specific test referencing it (e.g. a moisture content test record)
type LabTest struct {
	ID          int        `json:"id"`
	Name        NullString `json:"name"`
	Type        string     `json:"test_type" db:"type"`
	Sample      int        `json:"sample" db:"sample"`
	StartDate   NullDate   `json:"start_date" db:"start_date"`
	EndDate     NullDate   `json:"end_date" db:"end_date"`
	PerformedBy NullString `json:"performed_by" db:"performed_by"`
	CheckedDate NullDate   `json:"checked_date" db:"checked_date"`
	CheckedBy   NullString `json:"checked_by" db:"checked_by"`
}

// LabTestResponse is the data returned by the API containing info about a lab test.
type LabTestResponse struct {
	ID           int        `json:"id"`
	Name         NullString `json:"name"`
	Type         string     `json:"test_type" db:"type"`
	Sample       int        `json:"sample" db:"sample"`
	Borehole     int        `json:"borehole" db:"borehole"`
	BoreholeName string     `json:"borehole_name" db:"borehole_name"`
	StartDate    NullDate   `json:"start_date" db:"start_date"`
	EndDate      NullDate   `json:"end_date" db:"end_date"`
	PerformedBy  NullString `json:"performed_by" db:"performed_by"`
	CheckedDate  NullDate   `json:"checked_date" db:"checked_date"`
	CheckedBy    NullString `json:"checked_by" db:"checked_by"`
	SampleName   string     `json:"sample_name" db:"sample_name"`
}

// MoistureTestRequest is the data needed to create a moisture test record
type MoistureTestRequest struct {
	StartDate      NullDate   `json:"start_date" db:"start_date"`
	EndDate        NullDate   `json:"end_date" db:"end_date"`
	PerformedBy    NullString `json:"performed_by" db:"performed_by"`
	CheckedDate    NullDate   `json:"checked_date" db:"checked_date"`
	CheckedBy      NullString `json:"checked_by" db:"checked_by"`
	TareMass       *float64   `json:"tare_mass,string" db:"tare_mass"`
	SamplePlusTare *float64   `json:"sample_plus_tare,string" db:"sample_plus_tare"`
	DryPlusTare    *float64   `json:"dry_plus_tare,string" db:"dry_plus_tare"`
}

// MoistureTestResponse contains all the data relating to a moisture content test
type MoistureTestResponse struct {
	// ID is the ID of the lab_test and moisture_test (1 : 0..1)
	ID             int        `json:"id"`
	Name           NullString `json:"name"`
	Type           string     `json:"test_type" db:"type"`
	Sample         int        `json:"sample" db:"sample"`
	Borehole       int        `json:"borehole" db:"borehole"`
	BoreholeName   string     `json:"borehole_name" db:"borehole_name"`
	StartDate      NullDate   `json:"start_date" db:"start_date"`
	EndDate        NullDate   `json:"end_date" db:"end_date"`
	PerformedBy    NullString `json:"performed_by" db:"performed_by"`
	CheckedDate    NullDate   `json:"checked_date" db:"checked_date"`
	CheckedBy      NullString `json:"checked_by" db:"checked_by"`
	SampleName     string     `json:"sample_name" db:"sample_name"`
	TareMass       *float64   `json:"tare_mass" db:"tare_mass"`
	SamplePlusTare *float64   `json:"sample_plus_tare" db:"sample_plus_tare"`
	DryPlusTare    *float64   `json:"dry_plus_tare" db:"dry_plus_tare"`
}

// GSATestRequest is the set of data that the client provides to
// start or update a grain size analysis test record
type GSATestRequest struct {
	StartDate      NullDate         `json:"start_date" db:"start_date"`
	EndDate        NullDate         `json:"end_date" db:"end_date"`
	PerformedBy    NullString       `json:"performed_by" db:"performed_by"`
	CheckedDate    NullDate         `json:"checked_date" db:"checked_date"`
	CheckedBy      NullString       `json:"checked_by" db:"checked_by"`
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
	Name           NullString         `json:"name"`
	Type           string             `json:"test_type" db:"type"`
	Sample         int                `json:"sample" db:"sample"`
	SampleName     string             `json:"sample_name" db:"sample_name"`
	Borehole       int                `json:"borehole" db:"borehole"`
	BoreholeName   string             `json:"borehole_name" db:"borehole_name"`
	StartDate      NullDate           `json:"start_date" db:"start_date"`
	EndDate        NullDate           `json:"end_date" db:"end_date"`
	PerformedBy    NullString         `json:"performed_by" db:"performed_by"`
	CheckedDate    NullDate           `json:"checked_date" db:"checked_date"`
	CheckedBy      NullString         `json:"checked_by" db:"checked_by"`
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

// labTestCtxMiddleware is used by lab test routes that have a test id in the URI
func (s *server) labTestCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		labTestID, err := strconv.Atoi(chi.URLParam(r, "labTestID"))
		if err != nil {
			log.Println("labTestID not supplied")
			http.Error(w, http.StatusText(404), 404)
			return
		}

		labTest, err := s.datastore.RetrieveLabTest(labTestID)
		if err != nil {
			log.Println("labTest was not found in DB", err)
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), labTestCtx, labTest)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// func (s *server) gsaTestCtxMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		labTest, ok := ctx.Value(labTestCtx).(LabTestResponse)
// 		if !ok {
// 			http.Error(w, http.StatusText(422), 422)
// 			return
// 		}

// 		gsaTest, err := s.datastore.RetrieveSi

// 	})
// }

func (s *server) labTestOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

func (s *server) singleLabTestOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, DELETE, PUT, OPTIONS")
	return
}

// listLabTests returns lab tests from a project.
// the project must be passed in the request context with the contextKey "projectCtx"
func (s *server) listLabTestsByProject(w http.ResponseWriter, req *http.Request) {
	var err error
	ctx := req.Context()
	project, ok := ctx.Value(projects_v1.ProjectCtx).(projects_v1.Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	var boreholeID int
	borehole := req.FormValue("borehole")
	if borehole != "" {
		boreholeID, err = strconv.Atoi(borehole)
		if err != nil {
			// if borehole can't be converted to an int, make sure boreholeID is zero.
			// this ignores the ?borehole query if it's not a valid integer.
			boreholeID = 0
		}
	}

	labTests, err := s.datastore.ListLabTestsByProject(project.ID, boreholeID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, req, labTests)
}

func (s *server) createLabTest(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	labTest := LabTest{}
	err := decoder.Decode(&labTest)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	newTest, err := s.datastore.CreateLabTest(labTest)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newTest)
}

func (s *server) createMoistureTest(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	moistureTest := MoistureTestRequest{}
	err := decoder.Decode(&moistureTest)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	// Get test base object from context
	ctx := req.Context()
	labTest, ok := ctx.Value(labTestCtx).(LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// create the moisture content record
	newMoistureTest, err := s.datastore.CreateMoistureTest(moistureTest, labTest.ID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newMoistureTest)
}

func (s *server) retrieveMoistureTest(w http.ResponseWriter, req *http.Request) {
	labTestID, err := strconv.Atoi(chi.URLParam(req, "labTestID"))
	if err != nil {
		log.Println("test id not supplied")
		http.Error(w, http.StatusText(404), 404)
		return
	}

	moistureTest, err := s.datastore.RetrieveMoistureTest(labTestID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(w, req, moistureTest)
}

func (s *server) retrieveGSATest(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	labTest, ok := ctx.Value(labTestCtx).(LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	sieveTest, err := s.datastore.RetrieveSieveTest(labTest.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(w, req, sieveTest)
}

// putLabTest allows updating a lab test with a PUT request
// it calls a specific handler based on the type of lab test
// already registered in the database.
func (s *server) putLabTest(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	labTest, ok := ctx.Value(labTestCtx).(LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	switch labTest.Type {
	case "moisture_content":
		s.putMoistureTest(w, req)
	case "grain_size_analysis":
		s.putGSATest(w, req)
	default:
		// default handler for when no test type was provided
		// note: this might be heading towards "unreachable" and due for refactor.
		// all lab tests should have a corresponding child record (one of the above cases should match)

		decoder := json.NewDecoder(req.Body)
		labTestReq := LabTest{}
		err := decoder.Decode(&labTestReq)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		updatedTestData := LabTest{
			ID:          labTest.ID,
			Name:        labTestReq.Name,
			StartDate:   labTestReq.StartDate,
			EndDate:     labTestReq.EndDate,
			Type:        labTestReq.Type,
			PerformedBy: labTestReq.PerformedBy,
			Sample:      labTestReq.Sample,
			CheckedBy:   labTestReq.CheckedBy,
			CheckedDate: labTestReq.CheckedDate,
		}

		updatedLabTest, err := s.datastore.UpdateLabTest(updatedTestData)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		render.Status(req, http.StatusOK)
		render.JSON(w, req, updatedLabTest)
	}
}

func (s *server) putMoistureTest(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	moistureTestReq := MoistureTestRequest{}
	err := decoder.Decode(&moistureTestReq)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	ctx := req.Context()
	labTest, ok := ctx.Value(labTestCtx).(LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	updated, err := s.datastore.UpdateMoistureTest(moistureTestReq, labTest.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(w, req, updated)
}

func (s *server) putGSATest(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	gsaTestReq := GSATestRequest{}
	err := decoder.Decode(&gsaTestReq)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	ctx := req.Context()
	labTest, ok := ctx.Value(labTestCtx).(LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	updated, err := s.datastore.UpdateGSATest(gsaTestReq, labTest.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(w, req, updated)
}

// deleteLabTest asks the datastore to delete a test record
func (s *server) deleteLabTest(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	labTest, ok := ctx.Value(labTestCtx).(LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	err := s.datastore.DeleteLabTest(labTest.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.NoContent(w, req)
	return
}

func (s *server) addSieveToGSA(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	labTest, ok := ctx.Value(labTestCtx).(LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	if labTest.Type != "grain_size_analysis" {
		http.Error(w, "can only add sieves to a grain size analysis test", 400)
		return
	}

	decoder := json.NewDecoder(req.Body)
	sieveRequest := GSADataRequest{}
	err := decoder.Decode(&sieveRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	created, err := s.datastore.AddSieve(sieveRequest, labTest.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, created)
}

func (s *server) putSieve(w http.ResponseWriter, req *http.Request) {
	sieveID, err := strconv.Atoi(chi.URLParam(req, "sieveID"))
	if err != nil {
		log.Println("sieve id not supplied")
		http.Error(w, http.StatusText(404), 404)
		return
	}

	ctx := req.Context()
	labTest, ok := ctx.Value(labTestCtx).(LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	if labTest.Type != "grain_size_analysis" {
		http.Error(w, "can only modify sieves for a grain size analysis test", 400)
		return
	}

	decoder := json.NewDecoder(req.Body)
	sieveRequest := GSADataRequest{}
	err = decoder.Decode(&sieveRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	modified, err := s.datastore.UpdateSieve(sieveRequest, labTest.ID, sieveID)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	render.Status(req, http.StatusOK)
	render.JSON(w, req, modified)
}

func (s *server) deleteSieve(w http.ResponseWriter, req *http.Request) {
	sieveID, err := strconv.Atoi(chi.URLParam(req, "sieveID"))
	if err != nil {
		log.Println("sieve id not supplied")
		http.Error(w, http.StatusText(404), 404)
		return
	}

	err = s.datastore.DeleteSieve(sieveID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.NoContent(w, req)
	return
}

func (s *server) retrieveLabTest(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	labTest, ok := ctx.Value(labTestCtx).(LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	switch labTest.Type {
	case "grain_size_analysis":
		s.retrieveGSATest(w, req)
	case "moisture_content":
		s.retrieveMoistureTest(w, req)
	default:
		// default response (just a basic lab test response)
		// this code ideally should not be reachable, but this
		// is a safe fallback.
		render.Status(req, http.StatusOK)
		render.JSON(w, req, labTest)
	}

}
