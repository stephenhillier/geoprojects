package field

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

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

func (s *App) programOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}
