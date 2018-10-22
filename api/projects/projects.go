package projects

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// Project represents an engineering project. It holds files and data associated with a single project
type Project struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Location      string `json:"location"`
	BoreholeCount int    `json:"borehole_count" db:"borehole_count"`
}

// listProjects returns a list of all project records
func (s *App) listProjects(w http.ResponseWriter, req *http.Request) {

	projects, err := s.repo.AllProjects()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, req, projects)
}

// createProject handles a post request to the projects endpoint and
// creates a new project record.
// Requires details about the new project in the request body.
func (s *App) createProject(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// take input from POST request and store in a new Project type
	project := Project{}
	err = decoder.Decode(&project, req.PostForm)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// create the new project record in db
	record, err := s.repo.CreateProject(project)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// return the new project record (including its id)
	render.Status(req, http.StatusCreated)
	render.JSON(w, req, record)
}

// projectOpts serves a response to an OPTIONS request with allowed methods
func (s *App) projectOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

// singleProjectOptions serves a response to an OPTIONS request with allowed methods
func (s *App) singleProjectOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, OPTIONS, DELETE")
	return
}

// projectDetail retrieves one project record from database
func (s *App) projectDetail(w http.ResponseWriter, req *http.Request) {
	projectID, err := strconv.Atoi(chi.URLParam(req, "projectID"))
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	project, err := s.repo.RetrieveProject(projectID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	render.JSON(w, req, project)
}

// deleteProject sets a project to be expired at the current time
func (s *App) deleteProject(w http.ResponseWriter, req *http.Request) {
	// get projectID from URL
	projectID, err := strconv.Atoi(chi.URLParam(req, "projectID"))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}

	// check if project exists and/or not already expired
	_, err = s.repo.RetrieveProject(projectID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	// if the project exists, go ahead and delete it
	err = s.repo.DeleteProject(projectID)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
