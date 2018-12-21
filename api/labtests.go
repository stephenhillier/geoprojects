package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

func (s *server) labTestOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

func (s *server) singleLabTestOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "DELETE, PUT, OPTIONS")
	return
}

// listLabTests returns lab tests from a project.
// the project must be passed in the request context with the contextKey "projectCtx"
func (s *server) listLabTestsByProject(w http.ResponseWriter, req *http.Request) {
	var err error
	ctx := req.Context()
	project, ok := ctx.Value(projectCtx).(Project)
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
		log.Println("test was not found in DB", err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(w, req, moistureTest)
}

// putLabTest allows updating a lab test with a PUT request
func (s *server) putLabTest(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	labTestReq := LabTest{}
	err := decoder.Decode(&labTestReq)
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
