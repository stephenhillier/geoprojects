package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stephenhillier/geoprojects/earthworks"
)

// PaginatedProjectResponse is a paginated API response containing a count of all projects
// and the current page of projects
type PaginatedProjectResponse struct {
	Count   int                  `json:"count"`
	Results []earthworks.Project `json:"results"`
}

// List returns a list of all project records
func (svc *ProjectSvc) List(w http.ResponseWriter, req *http.Request) {

	limit, err := strconv.Atoi(req.FormValue("limit"))
	if err != nil || limit > svc.Settings.MaxPageLimit || limit < 0 {
		limit = svc.Settings.DefaultPageLimit
	}

	offset, err := strconv.Atoi(req.FormValue("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	projectName := req.FormValue("project_name")
	projectNumber := req.FormValue("project_number")

	// generic search that searches in both name and number (and other columns, potentially...)
	projectSearch := req.FormValue("search")

	projects, err := svc.Repo.AllProjects(projectName, projectNumber, projectSearch)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, req, projects)
}

// Create handles a post request to the projects endpoint and
// creates a new project record.
// Requires details about the new project in the request body.
func (svc *ProjectSvc) Create(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	// take input from POST request and store in a new Project type
	project := earthworks.ProjectRequest{}
	err := decoder.Decode(&project)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// create the new project record in db
	record, err := svc.Repo.CreateProject(project)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// return the new project record (including its id)
	render.Status(req, http.StatusCreated)
	render.JSON(w, req, record)
}

// Options serves a response to an OPTIONS request with allowed methods
func (svc *ProjectSvc) Options(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

// ProjectDetailOptions serves a response to an OPTIONS request with allowed methods
func (svc *ProjectSvc) ProjectDetailOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, OPTIONS, DELETE")
	return
}

// Retrieve retrieves one project record from database
func (svc *ProjectSvc) Retrieve(w http.ResponseWriter, req *http.Request) {
	// get project from request context
	project, err := getProjectContext(req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(422), 422)
		return
	}

	render.JSON(w, req, project)
}

// Update updates a project using a ProjectRequest
func (svc *ProjectSvc) Update(w http.ResponseWriter, req *http.Request) {
	// get project from request context
	project, err := getProjectContext(req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(422), 422)
		return
	}

	decoder := json.NewDecoder(req.Body)
	pReq := earthworks.ProjectRequest{}
	err = decoder.Decode(&pReq)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	updated, err := svc.Repo.UpdateProject(project.ID, pReq)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	render.JSON(w, req, updated)
}

// Delete deletes a project
func (svc *ProjectSvc) Delete(w http.ResponseWriter, req *http.Request) {

	// get project from request context
	project, err := getProjectContext(req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// if the project exists, go ahead and delete it
	err = svc.Repo.DeleteProject(project.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ProjectCtxMiddleware is used by Project routes that have a projectID in the URL path.
// it finds the specified project (returning 404 if the project is not found) and adds it
// to the request context.
func (svc *ProjectSvc) ProjectCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		project, err := svc.Repo.RetrieveProject(projectID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), earthworks.ContextKey{Name: "ProjectContext"}, project)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getProjectContext(r *http.Request) (earthworks.Project, error) {
	ctx := r.Context()
	project, ok := ctx.Value(earthworks.ContextKey{Name: "ProjectContext"}).(earthworks.Project)
	if !ok {
		return project, errors.New("error getting project from request context")
	}
	return project, nil
}
