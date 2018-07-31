package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/stephenhillier/geoprojects/backend/models"
)

var decoder = schema.NewDecoder()

// ProjectsIndex returns a list of all project records
func (api *Server) ProjectsIndex(w http.ResponseWriter, req *http.Request) {

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

// ProjectPost handles a post request to the projects endpoint and
// creates a new project record.
// Requires details about the new project in the request body.
func (api *Server) ProjectPost(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	if req.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, http.StatusText(405), 405)
	}

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
	w.Write(response)
}
