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

// BoreholeCreateRequest is the data a user should submit to create a borehole.
// A borehole can either be associated with an existing datapoint, or if a location
// is supplied, a datapoint will be created.
type BoreholeCreateRequest struct {
	Project   int64      `json:"project"`
	Program   NullInt64  `json:"program"`
	Datapoint NullInt64  `json:"datapoint"`
	Name      string     `json:"name"`
	StartDate NullDate   `json:"start_date" db:"start_date" schema:"start_date"`
	EndDate   NullDate   `json:"end_date" db:"end_date" schema:"end_date"`
	FieldEng  string     `json:"field_eng" db:"field_eng" schema:"field_eng"`
	Location  [2]float64 `json:"location"`
}

// BoreholeResponse is the data returned by the API after receiving a request for
// a borehole's details
// the FieldEng field is a string (users.username) instead of a primary key reference.
type BoreholeResponse struct {
	ID        int64     `json:"id"`
	Project   NullInt64 `json:"project"`
	Program   NullInt64 `json:"program"`
	Datapoint NullInt64 `json:"datapoint"`
	Name      string    `json:"name"`
	StartDate NullDate  `json:"start_date" db:"start_date"`
	EndDate   NullDate  `json:"end_date" db:"end_date"`
	FieldEng  string    `json:"field_eng" db:"field_eng"`
}

// PaginatedBoreholeResponse contains a count of all borehole records and paginated results from the database
type PaginatedBoreholeResponse struct {
	Count   int64               `json:"count"`
	Results []*BoreholeResponse `json:"results"`
}

// Borehole is drilled geotechnical test hole located at a Datapoint.
// There may be a number of samples/observations associated with one borehole.
type Borehole struct {
	ID        int64     `json:"id"`
	Project   NullInt64 `json:"project"`
	Program   NullInt64 `json:"program"`
	Datapoint int64     `json:"datapoint"`
	Name      string    `json:"name"`
	StartDate NullDate  `json:"start_date" db:"start_date"`
	EndDate   NullDate  `json:"end_date" db:"end_date"`
	FieldEng  string    `json:"field_eng" db:"field_eng"`
}

// boreholeOptions responds to OPTIONS requests (and pre-flight requests)
func (s *server) boreholeOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

// listBoreholes returns a list of boreholes for a specified project
// the project should be specified in a URL query string e.g. api/v1/boreholes?project=1
func (s *server) listBoreholes(w http.ResponseWriter, req *http.Request) {

	project := req.FormValue("project")
	limit, err := strconv.Atoi(req.FormValue("limit"))
	if err != nil || limit > s.config.maxPageLimit || limit < 0 {
		limit = s.config.defaultPageLimit
	}

	offset, err := strconv.Atoi(req.FormValue("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	var projectID int

	// if a project was supplied in querystring, set projectID so that the db query can
	// list boreholes by project
	if project != "" {
		projectID, err = strconv.Atoi(project)
		if err != nil {
			// if project can't be converted to an int, make sure projectID is zero.
			// this ignores the ?project query if it's not a valid integer.
			projectID = 0
		}
	}

	boreholes, count, err := s.datastore.ListBoreholes(projectID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	page := PaginatedBoreholeResponse{
		Count:   count,
		Results: boreholes,
	}

	render.JSON(w, req, page)
}

// createBorehole creates a new borehole
func (s *server) createBorehole(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	borehole := BoreholeCreateRequest{}
	err := decoder.Decode(&borehole)
	log.Println(borehole)
	if err != nil {
		log.Println(err)

		http.Error(w, err.Error(), 400)
		return
	}

	newBorehole, err := s.datastore.CreateBorehole(borehole)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newBorehole)
}

// func getBorehole(w http.ResponseWriter, req *http.Request) {

// }

// boreholeCtxMiddleware is used by borehole routes that have a boreholeID in the URL path.
// it passes borehole information into the request context
func (s *server) boreholeCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		boreholeID, err := strconv.Atoi(chi.URLParam(r, "boreholeID"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		borehole, err := s.datastore.GetBorehole(boreholeID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		log.Println(borehole)

		ctx := context.WithValue(r.Context(), boreholeCtx, borehole)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
