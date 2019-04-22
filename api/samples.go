package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	boreholev1 "github.com/stephenhillier/geoprojects/api/boreholes/model"
	projects_v1 "github.com/stephenhillier/geoprojects/api/projects/model"
)

// Sample is a soil layer/stratum and contains information such as description and depth of the layer
type Sample struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Borehole     int64   `json:"borehole"`
	BoreholeName string  `json:"borehole_name" db:"borehole_name"`
	Start        float64 `json:"start" db:"start_depth"`
	End          float64 `json:"end" db:"end_depth"`
	Description  string  `json:"description"`
	USCS         string  `json:"uscs" db:"uscs"`
	Tests        int     `json:"test_count"`
}

// SampleRequest is a struct containing fields required to create a new sample
type SampleRequest struct {
	Name        string  `json:"name"`
	Borehole    int64   `json:"borehole,string"`
	Start       float64 `json:"start,string"`
	End         float64 `json:"end,string"`
	Description string  `json:"description"`
	USCS        string  `json:"uscs" db:"uscs"`
}

func (s *server) sampleOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, PUT, DELETE, OPTIONS")
	return
}

// listSamplesByBorehole returns samples associated with a specified borehole.
// the borehole must be passed in the request context with the contextKey "boreholeCtx"
func (s *server) listSamplesByBorehole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	borehole, ok := ctx.Value(boreholev1.BoreholeCtx).(boreholev1.BoreholeResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	samples, err := s.datastore.ListSamplesByBorehole(borehole.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, req, samples)
}

func (s *server) listSamplesByProject(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	project, ok := ctx.Value(projects_v1.ProjectCtx).(projects_v1.Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	samples, err := s.datastore.ListSamplesByProject(project.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, req, samples)
}

// createSample parses a sample request and passes a sample object to the datastore to be created.
// if successful, it returns a new Sample object with the stored data
func (s *server) createSample(w http.ResponseWriter, req *http.Request) {

	// get borehole context (this is the borehole indicated by the borehole number in the URI)
	ctx := req.Context()
	borehole, ok := ctx.Value(boreholev1.BoreholeCtx).(boreholev1.BoreholeResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	decoder := json.NewDecoder(req.Body)
	sampleReq := SampleRequest{}
	err := decoder.Decode(&sampleReq)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	sample := Sample{
		Borehole:    borehole.ID,
		Name:        sampleReq.Name,
		Start:       sampleReq.Start,
		End:         sampleReq.End,
		Description: sampleReq.Description,
		USCS:        sampleReq.USCS,
	}

	newSample, err := s.datastore.CreateSample(sample)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newSample)
}

// putSample allows updating a sample record by making a PUT request to the sample's endpoint
func (s *server) putSample(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	sampleReq := SampleRequest{}
	err := decoder.Decode(&sampleReq)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	ctx := req.Context()
	sample, ok := ctx.Value(sampleCtx).(Sample)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	updatedSampleData := Sample{
		ID:          sample.ID,
		Name:        sampleReq.Name,
		Borehole:    sampleReq.Borehole,
		Start:       sampleReq.Start,
		End:         sampleReq.End,
		Description: sampleReq.Description,
		USCS:        sampleReq.USCS,
	}

	updatedSample, err := s.datastore.UpdateSample(updatedSampleData)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusOK)
	render.JSON(w, req, updatedSample)
}

// deleteSample asks the datastore to delete a given sample record
func (s *server) deleteSample(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	sample, ok := ctx.Value(sampleCtx).(Sample)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	err := s.datastore.DeleteSample(sample.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.NoContent(w, req)
	return
}

func (s *server) retrieveSample(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	sample, ok := ctx.Value(sampleCtx).(Sample)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(w, req, sample)
}

// sampleCtxMiddleware is used by sample routes that have a sampleID in the URL path.
// it finds the specified sample (returning 404 if the sample is not found) and adds it
// to the request context.
func (s *server) sampleCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sampleID, err := strconv.Atoi(chi.URLParam(r, "sampleID"))
		if err != nil {
			log.Println("sampleID not supplied")
			http.Error(w, http.StatusText(404), 404)
			return
		}

		sample, err := s.datastore.RetrieveSample(sampleID)
		if err != nil {
			log.Println("sample was not found in DB")
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), sampleCtx, sample)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
