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

// SampleOptions responds, with allowed methods, to a preflight options request
func (svc *BoreholeSvc) SampleOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, PUT, DELETE, OPTIONS")
	return
}

// ListSamplesByBorehole returns samples associated with a specified borehole.
// the borehole must be passed in the request context with the contextKey "boreholeCtx"
func (svc *BoreholeSvc) ListSamplesByBorehole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	borehole, ok := ctx.Value(earthworks.BoreholeCtx).(earthworks.BoreholeResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	samples, err := svc.Repo.ListSamplesByBorehole(borehole.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, req, samples)
}

// ListSamplesByProject returns samples from a project specified in the URL path
func (svc *BoreholeSvc) ListSamplesByProject(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	project, ok := ctx.Value(earthworks.ContextKey{Name: "ProjectContext"}).(earthworks.Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	samples, err := svc.Repo.ListSamplesByProject(project.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, req, samples)
}

// CreateSample parses a sample request and passes a sample object to the datastore to be created.
// if successful, it returns a new Sample object with the stored data
func (svc *BoreholeSvc) CreateSample(w http.ResponseWriter, req *http.Request) {

	// get borehole context (this is the borehole indicated by the borehole number in the URI)
	ctx := req.Context()
	borehole, ok := ctx.Value(earthworks.BoreholeCtx).(earthworks.BoreholeResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	decoder := json.NewDecoder(req.Body)
	sampleReq := earthworks.SampleRequest{}
	err := decoder.Decode(&sampleReq)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	sample := earthworks.Sample{
		Borehole:    borehole.ID,
		Name:        sampleReq.Name,
		Start:       sampleReq.Start,
		End:         sampleReq.End,
		Description: sampleReq.Description,
		USCS:        sampleReq.USCS,
	}

	newSample, err := svc.Repo.CreateSample(sample)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newSample)
}

// PutSample allows updating a sample record by making a PUT request to the sample's endpoint
func (svc *BoreholeSvc) PutSample(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	sampleReq := earthworks.SampleRequest{}
	err := decoder.Decode(&sampleReq)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	ctx := req.Context()
	sample, ok := ctx.Value(earthworks.SampleCtx).(earthworks.Sample)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	updatedSampleData := earthworks.Sample{
		ID:          sample.ID,
		Name:        sampleReq.Name,
		Borehole:    sampleReq.Borehole,
		Start:       sampleReq.Start,
		End:         sampleReq.End,
		Description: sampleReq.Description,
		USCS:        sampleReq.USCS,
	}

	updatedSample, err := svc.Repo.UpdateSample(updatedSampleData)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusOK)
	render.JSON(w, req, updatedSample)
}

// DeleteSample asks the datastore to delete a given sample record
func (svc *BoreholeSvc) DeleteSample(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	sample, ok := ctx.Value(earthworks.SampleCtx).(earthworks.Sample)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	err := svc.Repo.DeleteSample(sample.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.NoContent(w, req)
	return
}

// RetrieveSample responds to a GET request and returns the sample
// specified by the URL path
func (svc *BoreholeSvc) RetrieveSample(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	sample, ok := ctx.Value(earthworks.SampleCtx).(earthworks.Sample)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(w, req, sample)
}

// SampleCtxMiddleware is used by sample routes that have a sampleID in the URL path.
// it finds the specified sample (returning 404 if the sample is not found) and adds it
// to the request context.
func (svc *BoreholeSvc) SampleCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sampleID, err := strconv.Atoi(chi.URLParam(r, "sampleID"))
		if err != nil {
			log.Println("sampleID not supplied")
			http.Error(w, http.StatusText(404), 404)
			return
		}

		sample, err := svc.Repo.RetrieveSample(sampleID)
		if err != nil {
			log.Println("sample was not found in DB")
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), earthworks.SampleCtx, sample)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
