package projects

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// App represents an HTTP web application with a datastore, handlers and routes.
// Routes can be passed into a chi.Router Route() to provide an
// access point to the handlers in this app.
type App struct {
	repo   Repository
	Routes func(r chi.Router)
}

// listProjects returns a list of all project records
func (s *App) listProjects(w http.ResponseWriter, req *http.Request) {

	if req.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	projects, err := s.repo.AllProjects()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	response, err := json.Marshal(projects)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// createProject handles a post request to the projects endpoint and
// creates a new project record.
// Requires details about the new project in the request body.
func (s *App) createProject(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	if req.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, http.StatusText(405), 405)
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
	response, err := json.Marshal(record)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// projectOpts serves a response to an OPTIONS request with allowed methods
func (s *App) projectOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

func (s *App) singleProjectOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, OPTIONS")
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

	response, err := json.Marshal(project)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
