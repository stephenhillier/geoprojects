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

// PaginatedBoreholeResponse contains a count of all borehole records and paginated results from the database
type PaginatedBoreholeResponse struct {
	Count   int64                          `json:"count"`
	Results []*earthworks.BoreholeResponse `json:"results"`
}

// Options responds to OPTIONS requests (and pre-flight requests)
func (svc *BoreholeSvc) Options(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

// List returns a list of boreholes for a specified project
// the project should be specified in a URL query string e.g. api/v1/boreholes?project=1
func (svc *BoreholeSvc) List(w http.ResponseWriter, req *http.Request) {

	project := req.FormValue("project")
	limit, err := strconv.Atoi(req.FormValue("limit"))
	if err != nil || limit > svc.Settings.MaxPageLimit || limit < 0 {
		limit = svc.Settings.DefaultPageLimit
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

	boreholes, count, err := svc.Repo.ListBoreholes(projectID, limit, offset)
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
func (svc *BoreholeSvc) Create(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	borehole := earthworks.BoreholeCreateRequest{}
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

	newBorehole, err := svc.Repo.CreateBorehole(borehole, projectID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newBorehole)
}

// Get retrieves a borehole
func (svc *BoreholeSvc) Get(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	borehole, ok := ctx.Value(earthworks.BoreholeCtx).(earthworks.BoreholeResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// TODO
	// strataCount, err := svc.Repo.CountStrataForBorehole(borehole.ID)
	// if err != nil {
	// 	log.Println("error fetching strata count:", err)
	// }

	// borehole.StrataCount = strataCount

	render.JSON(w, req, borehole)
}

// Delete asks the datastore to delete a given borehole record
func (svc *BoreholeSvc) Delete(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	borehole, ok := ctx.Value(earthworks.BoreholeCtx).(earthworks.BoreholeResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	err := svc.Repo.DeleteBorehole(borehole.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.NoContent(w, req)
	return
}

// BoreholeCtxMiddleware is used by borehole routes that have a boreholeID in the URL path.
// it passes borehole information into the request context
func (svc *BoreholeSvc) BoreholeCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		boreholeID, err := strconv.Atoi(chi.URLParam(r, "boreholeID"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		borehole, err := svc.Repo.GetBorehole(boreholeID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), earthworks.BoreholeCtx, borehole)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
