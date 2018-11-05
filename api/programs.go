package main

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func (s *server) programOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

func (s *server) listPrograms(w http.ResponseWriter, req *http.Request) {
	programs, err := s.datastore.ListPrograms()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, req, programs)
}

func (s *server) createProgram(w http.ResponseWriter, req *http.Request) {
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
	newRecord, err := s.datastore.CreateProgram(program)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newRecord)
}
