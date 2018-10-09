package projects

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

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

func (s *App) singleProjectOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, OPTIONS")
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
