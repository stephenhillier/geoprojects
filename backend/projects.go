package main

import (
	"encoding/json"
	"net/http"
)

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

	w.Write(response)
}
