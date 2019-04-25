package server

import (
	"github.com/go-chi/chi"
	"github.com/stephenhillier/geoprojects/api/boreholes"
	"github.com/stephenhillier/geoprojects/api/db"
	"github.com/stephenhillier/geoprojects/api/projects"
	"github.com/stephenhillier/geoprojects/api/server/config"
)

// Server represents the server environment (db and router)
type Server struct {
	router    chi.Router
	datastore db.Datastore
	config    config.Config
	handlers  handlers
}

type handlers struct {
	Projects  *projects.ProjectSvc
	Boreholes *boreholes.BoreholeSvc
}
