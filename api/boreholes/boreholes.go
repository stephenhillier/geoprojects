package boreholes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	boreholev1 "github.com/stephenhillier/geoprojects/api/boreholes/model"
	"github.com/stephenhillier/geoprojects/api/boreholes/repository"
	"github.com/stephenhillier/geoprojects/api/db"
	"github.com/stephenhillier/geoprojects/api/server/config"
)

// PaginatedBoreholeResponse contains a count of all borehole records and paginated results from the database
type PaginatedBoreholeResponse struct {
	Count   int64                          `json:"count"`
	Results []*boreholev1.BoreholeResponse `json:"results"`
}

// BoreholeSvc is a service that provides methods for working with boreholes
type BoreholeSvc struct {
	repo   repository.BoreholeRepository
	config *config.Config
}

// NewBoreholeSvc returns a BoreholeSvc with methods for working with boreholes
func NewBoreholeSvc(store *db.Datastore, config *config.Config) *BoreholeSvc {
	return &BoreholeSvc{
		config: config,
		repo:   repository.NewBoreholeRepo(store),
	}
}

// Options responds to OPTIONS requests (and pre-flight requests)
func (s *BoreholeSvc) Options(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

// List returns a list of boreholes for a specified project
// the project should be specified in a URL query string e.g. api/v1/boreholes?project=1
func (s *BoreholeSvc) List(w http.ResponseWriter, req *http.Request) {

	project := req.FormValue("project")
	limit, err := strconv.Atoi(req.FormValue("limit"))
	if err != nil || limit > s.config.MaxPageLimit || limit < 0 {
		limit = s.config.DefaultPageLimit
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

	boreholes, count, err := s.repo.ListBoreholes(projectID, limit, offset)
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

// Create creates a new borehole
func (s *BoreholeSvc) Create(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	borehole := boreholev1.BoreholeCreateRequest{}
	err := decoder.Decode(&borehole)
	if err != nil {
		log.Println(err)

		http.Error(w, err.Error(), 400)
		return
	}

	projectID, err := borehole.Project.Int64()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	newBorehole, err := s.repo.CreateBorehole(borehole, projectID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newBorehole)
}

// Get retrieves a borehole
func (s *BoreholeSvc) Get(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	borehole, ok := ctx.Value(boreholev1.BoreholeCtx).(boreholev1.BoreholeResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// TODO
	// strataCount, err := s.repo.CountStrataForBorehole(borehole.ID)
	// if err != nil {
	// 	log.Println("error fetching strata count:", err)
	// }

	// borehole.StrataCount = strataCount

	render.JSON(w, req, borehole)
}

// Delete asks the datastore to delete a given borehole record
func (s *BoreholeSvc) Delete(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	borehole, ok := ctx.Value(boreholev1.BoreholeCtx).(boreholev1.BoreholeResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	err := s.repo.DeleteBorehole(borehole.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.NoContent(w, req)
	return
}

// BoreholeCtxMiddleware is used by borehole routes that have a boreholeID in the URL path.
// it passes borehole information into the request context
func (s *BoreholeSvc) BoreholeCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		boreholeID, err := strconv.Atoi(chi.URLParam(r, "boreholeID"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		borehole, err := s.repo.GetBorehole(boreholeID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), boreholev1.BoreholeCtx, borehole)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
