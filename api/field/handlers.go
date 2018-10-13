package field

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func (s *App) programOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

func (s *App) boreholeOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

func (s *App) listPrograms(w http.ResponseWriter, req *http.Request) {
	programs, err := s.programs.ListPrograms()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, req, programs)
}

func (s *App) createProgram(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	program := ProgramCreateRequest{}
	err = decoder.Decode(&program, req.PostForm)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// create the new program record
	newRecord, err := s.programs.CreateProgram(program)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newRecord)
}

func (s *App) listBoreholes(w http.ResponseWriter, req *http.Request) {

	project := req.FormValue("project")

	var projectID int
	var err error

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

	boreholes, err := s.boreholes.ListBoreholes(projectID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	render.JSON(w, req, boreholes)
}

func (s *App) createBorehole(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	borehole := BoreholeCreateRequest{}
	err := decoder.Decode(&borehole)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	newBorehole, err := s.boreholes.CreateBorehole(borehole)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newBorehole)
}

// func getBorehole(w http.ResponseWriter, req *http.Request) {

// }
