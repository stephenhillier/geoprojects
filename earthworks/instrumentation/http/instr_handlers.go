package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/stephenhillier/geoprojects/earthworks"
)

// List displays all instrumentation for a given project
func (svc *InstrumentationSvc) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project, ok := ctx.Value(earthworks.ProjectCtx).(earthworks.Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	instruments, err := svc.Repo.ListInstruments(project.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	render.JSON(w, r, instruments)

}

// Create handles a POST request and creates a new instrument in the repo
func (svc *InstrumentationSvc) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	instrRequest := earthworks.InstrumentCreateRequest{}
	err := decoder.Decode(&instrRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	// get project context
	ctx := r.Context()
	project, ok := ctx.Value(earthworks.ProjectCtx).(earthworks.Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	newInstrument, err := svc.Repo.CreateInstrument(instrRequest, project.ID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, newInstrument)

}
