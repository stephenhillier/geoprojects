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

// Project represents an engineering project. It holds files and data associated with a single project
type Project struct {
	ID int `json:"id"`

	// Number refers to a Project Number.  This is a typical industry term but it often contains letters
	// (to indicate the region, department etc)
	Number           string        `json:"number"`
	Name             string        `json:"name"`
	Location         string        `json:"location"`
	Client           string        `json:"client"`
	PM               string        `json:"pm"`
	BoreholeCount    int           `json:"borehole_count" db:"borehole_count"`
	CentroidLocation PointLocation `json:"centroid" db:"centroid"`
	DefaultCoords    PointLocation `json:"default_coords" db:"default_coords"`
}

// ProjectRequest is the set of data required to accept a request for a new project
type ProjectRequest struct {
	Name          string     `json:"name"`
	Number        string     `json:"number"`
	Client        string     `json:"client"`
	PM            string     `json:"pm"`
	Location      string     `json:"location"`
	DefaultCoords [2]float64 `json:"default_coords"`
}

// PaginatedProjectResponse is a paginated API response containing a count of all projects
// and the current page of projects
type PaginatedProjectResponse struct {
	Count   int       `json:"count"`
	Results []Project `json:"results"`
}

// listProjects returns a list of all project records
func (s *server) listProjects(w http.ResponseWriter, req *http.Request) {

	limit, err := strconv.Atoi(req.FormValue("limit"))
	if err != nil || limit > s.config.maxPageLimit || limit < 0 {
		limit = s.config.defaultPageLimit
	}

	offset, err := strconv.Atoi(req.FormValue("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	projectName := req.FormValue("project_name")
	projectNumber := req.FormValue("project_number")

	// generic search that searches in both name and number (and other columns, potentially...)
	projectSearch := req.FormValue("search")

	projects, err := s.datastore.AllProjects(projectName, projectNumber, projectSearch)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, req, projects)
}

// createProject handles a post request to the projects endpoint and
// creates a new project record.
// Requires details about the new project in the request body.
func (s *server) createProject(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	// take input from POST request and store in a new Project type
	project := ProjectRequest{}
	err := decoder.Decode(&project)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// create the new project record in db
	record, err := s.datastore.CreateProject(project)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// return the new project record (including its id)
	render.Status(req, http.StatusCreated)
	render.JSON(w, req, record)
}

// projectOpts serves a response to an OPTIONS request with allowed methods
func (s *server) projectOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

// singleProjectOptions serves a response to an OPTIONS request with allowed methods
func (s *server) singleProjectOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, OPTIONS, DELETE")
	return
}

// projectDetail retrieves one project record from database
func (s *server) projectDetail(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	project, ok := ctx.Value(projectCtx).(Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	render.JSON(w, req, project)
}

// deleteProject deletes a project
func (s *server) deleteProject(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	project, ok := ctx.Value(projectCtx).(Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// if the project exists, go ahead and delete it
	err := s.datastore.DeleteProject(project.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// projectCtxMiddleware is used by Project routes that have a projectID in the URL path.
// it finds the specified project (returning 404 if the project is not found) and adds it
// to the request context.
func (s *server) projectCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		project, err := s.datastore.RetrieveProject(projectID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), projectCtx, project)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
