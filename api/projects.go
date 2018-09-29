package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/gorilla/schema"
	"github.com/stephenhillier/geoprojects/backend/models"
)

var decoder = schema.NewDecoder()

// listProjects returns a list of all project records
func (api *Server) listProjects(w http.ResponseWriter, req *http.Request) {

	if req.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	projects, err := api.db.AllProjects()
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
func (api *Server) createProject(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	if req.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, http.StatusText(405), 405)
	}

	// take input from POST request and store in a new Project type
	project := models.Project{}
	err = decoder.Decode(&project, req.PostForm)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// create the new project record in db
	record, err := api.db.CreateProject(project)
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
func (api *Server) projectOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

func (api *Server) singleProjectOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, OPTIONS")
}

// projectDetail retrieves one project record from database
func (api *Server) projectDetail(w http.ResponseWriter, req *http.Request) {
	projectID, err := strconv.Atoi(chi.URLParam(req, "projectID"))
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	project, err := api.db.RetrieveProject(projectID)
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
