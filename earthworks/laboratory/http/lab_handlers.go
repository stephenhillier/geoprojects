package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stephenhillier/geoprojects/earthworks"
)

// LabTestCtxMiddleware is used by lab test routes that have a test id in the URI
func (s *LabSvc) LabTestCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		labTestID, err := strconv.Atoi(chi.URLParam(r, "labTestID"))
		if err != nil {
			log.Println("labTestID not supplied")
			http.Error(w, http.StatusText(404), 404)
			return
		}

		labTest, err := s.Repo.RetrieveLabTest(labTestID)
		if err != nil {
			log.Println("labTest was not found in DB", err)
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), earthworks.LabTestCtx, labTest)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// func (s *LabSvc) gsaTestCtxMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		labTest, ok := ctx.Value(earthworks.LabTestCtx).(earthworks.LabTestResponse)
// 		if !ok {
// 			http.Error(w, http.StatusText(422), 422)
// 			return
// 		}

// 		gsaTest, err := s.Repo.RetrieveSi

// 	})
// }

// LabTestOptions handles an options request for lab test routes
func (s *LabSvc) LabTestOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

// SingleLabTestOptions handles an options request for single lab test routes
func (s *LabSvc) SingleLabTestOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, DELETE, PUT, OPTIONS")
	return
}

// ListLabTestsByProject returns lab tests from a project.
// the project must be passed in the request context with the contextKey "projectCtx"
func (s *LabSvc) ListLabTestsByProject(w http.ResponseWriter, req *http.Request) {
	var err error
	ctx := req.Context()
	project, ok := ctx.Value(earthworks.ProjectCtx).(earthworks.Project)
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

	labTests, err := s.Repo.ListLabTestsByProject(project.ID, boreholeID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, req, labTests)
}

// CreateLabTest handles a request to create a lab test
func (s *LabSvc) CreateLabTest(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	labTest := earthworks.LabTest{}
	err := decoder.Decode(&labTest)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	newTest, err := s.Repo.CreateLabTest(labTest)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newTest)
}

// CreateMoistureTest handles a request to create a moisture test
func (s *LabSvc) CreateMoistureTest(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	moistureTest := earthworks.MoistureTestRequest{}
	err := decoder.Decode(&moistureTest)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	// Get test base object from context
	ctx := req.Context()
	labTest, ok := ctx.Value(earthworks.LabTestCtx).(earthworks.LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// create the moisture content record
	newMoistureTest, err := s.Repo.CreateMoistureTest(moistureTest, labTest.ID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newMoistureTest)
}

// RetrieveMoistureTest handles a request for details of a moisture content test
func (s *LabSvc) RetrieveMoistureTest(w http.ResponseWriter, req *http.Request) {
	labTestID, err := strconv.Atoi(chi.URLParam(req, "labTestID"))
	if err != nil {
		log.Println("test id not supplied")
		http.Error(w, http.StatusText(404), 404)
		return
	}

	moistureTest, err := s.Repo.RetrieveMoistureTest(labTestID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(w, req, moistureTest)
}

// RetrieveGSATest handles a request for details of a grain size test
func (s *LabSvc) RetrieveGSATest(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	labTest, ok := ctx.Value(earthworks.LabTestCtx).(earthworks.LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	sieveTest, err := s.Repo.RetrieveSieveTest(labTest.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(w, req, sieveTest)
}

// PutLabTest allows updating a lab test with a PUT request
// it calls a specific handler based on the type of lab test
// already registered in the database.
func (s *LabSvc) PutLabTest(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	labTest, ok := ctx.Value(earthworks.LabTestCtx).(earthworks.LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	switch labTest.Type {
	case "moisture_content":
		s.PutMoistureTest(w, req)
	case "grain_size_analysis":
		s.PutGSATest(w, req)
	default:
		// default handler for when no test type was provided
		// note: this might be heading towards "unreachable" and due for refactor.
		// all lab tests should have a corresponding child record (one of the above cases should match)

		decoder := json.NewDecoder(req.Body)
		labTestReq := earthworks.LabTest{}
		err := decoder.Decode(&labTestReq)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		updatedTestData := earthworks.LabTest{
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

		updatedLabTest, err := s.Repo.UpdateLabTest(updatedTestData)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		render.Status(req, http.StatusOK)
		render.JSON(w, req, updatedLabTest)
	}
}

// PutMoistureTest modifies a moisture content test
func (s *LabSvc) PutMoistureTest(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	moistureTestReq := earthworks.MoistureTestRequest{}
	err := decoder.Decode(&moistureTestReq)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	ctx := req.Context()
	labTest, ok := ctx.Value(earthworks.LabTestCtx).(earthworks.LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	updated, err := s.Repo.UpdateMoistureTest(moistureTestReq, labTest.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(w, req, updated)
}

// PutGSATest modifies a grain size test
func (s *LabSvc) PutGSATest(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	gsaTestReq := earthworks.GSATestRequest{}
	err := decoder.Decode(&gsaTestReq)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	ctx := req.Context()
	labTest, ok := ctx.Value(earthworks.LabTestCtx).(earthworks.LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	updated, err := s.Repo.UpdateGSATest(gsaTestReq, labTest.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(w, req, updated)
}

// DeleteLabTest asks the datastore to delete a test record
func (s *LabSvc) DeleteLabTest(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	labTest, ok := ctx.Value(earthworks.LabTestCtx).(earthworks.LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	err := s.Repo.DeleteLabTest(labTest.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.NoContent(w, req)
	return
}

// AddSieveToGSA adds one sieve record to a grain size test
func (s *LabSvc) AddSieveToGSA(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	labTest, ok := ctx.Value(earthworks.LabTestCtx).(earthworks.LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	if labTest.Type != "grain_size_analysis" {
		http.Error(w, "can only add sieves to a grain size analysis test", 400)
		return
	}

	decoder := json.NewDecoder(req.Body)
	sieveRequest := earthworks.GSADataRequest{}
	err := decoder.Decode(&sieveRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	created, err := s.Repo.AddSieve(sieveRequest, labTest.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, created)
}

// PutSieve modifies a single sieve
func (s *LabSvc) PutSieve(w http.ResponseWriter, req *http.Request) {
	sieveID, err := strconv.Atoi(chi.URLParam(req, "sieveID"))
	if err != nil {
		log.Println("sieve id not supplied")
		http.Error(w, http.StatusText(404), 404)
		return
	}

	ctx := req.Context()
	labTest, ok := ctx.Value(earthworks.LabTestCtx).(earthworks.LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	if labTest.Type != "grain_size_analysis" {
		http.Error(w, "can only modify sieves for a grain size analysis test", 400)
		return
	}

	decoder := json.NewDecoder(req.Body)
	sieveRequest := earthworks.GSADataRequest{}
	err = decoder.Decode(&sieveRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	modified, err := s.Repo.UpdateSieve(sieveRequest, labTest.ID, sieveID)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	render.Status(req, http.StatusOK)
	render.JSON(w, req, modified)
}

// DeleteSieve deletes a single sieve from a grain size test
func (s *LabSvc) DeleteSieve(w http.ResponseWriter, req *http.Request) {
	sieveID, err := strconv.Atoi(chi.URLParam(req, "sieveID"))
	if err != nil {
		log.Println("sieve id not supplied")
		http.Error(w, http.StatusText(404), 404)
		return
	}

	err = s.Repo.DeleteSieve(sieveID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.NoContent(w, req)
	return
}

// RetrieveLabTest responds with the requested lab test record.
// the basic lab test data will already be on the request context since this handler
// uses the LabTestCtxMiddleware.
// The test-specific data (for moisture, grain size tests etc) will be
// retrieved based on the type of the test.
func (s *LabSvc) RetrieveLabTest(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	labTest, ok := ctx.Value(earthworks.LabTestCtx).(earthworks.LabTestResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	switch labTest.Type {
	case "grain_size_analysis":
		s.RetrieveGSATest(w, req)
	case "moisture_content":
		s.RetrieveMoistureTest(w, req)
	default:
		// default response (just a basic lab test response)
		// this code ideally should not be reachable, but this
		// is a safe fallback.
		render.Status(req, http.StatusOK)
		render.JSON(w, req, labTest)
	}

}
