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
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Location      string `json:"location"`
	BoreholeCount int    `json:"borehole_count" db:"borehole_count"`
}

// listProjects returns a list of all project records
func (s *server) listProjects(w http.ResponseWriter, req *http.Request) {

	projects, err := s.datastore.AllProjects()
	if err != nil {
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
	project := Project{}
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

// deleteProject sets a project to be expired at the current time
func (s *server) deleteProject(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	project, ok := ctx.Value(projectCtx).(*Project)
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

		log.Println(project)

		ctx := context.WithValue(r.Context(), projectCtx, project)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
