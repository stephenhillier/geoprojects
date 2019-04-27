package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/stephenhillier/geoprojects/earthworks"
	"github.com/stephenhillier/geoprojects/earthworks/db"
	projectsHandlers "github.com/stephenhillier/geoprojects/earthworks/projects/http"
)

// Service is a type that contains backing services and settings
// for the Earthworks backend app.
type Service struct {
	Settings earthworks.Settings
	Router   chi.Router
	Handlers Handlers
}

// Handlers contains http request handlers for each service
type Handlers struct {
	Projects *projectsHandlers.ProjectSvc
}

// NewEarthworksService returns an EarthworksService object that ties
// together the various services that form the backend Earthworks application
func NewEarthworksService() (Service, error) {

	svc := Service{}

	dbConfig := db.Config{
		Conn:   "postgres://127.0.0.1:5432/geo?sslmode=disable",
		Driver: "postgres",
	}
	db, err := db.NewDB(dbConfig)
	if err != nil {
		return Service{}, err
	}

	settings := earthworks.Settings{
		DefaultPageLimit: 10,
		MaxPageLimit:     100,
	}

	projects := projectsHandlers.NewProjectSvc(db, settings)
	// boreholes := boreholes.NewBoreholeSvc(store, cnf)

	log.Println(projects)

	r := chi.NewRouter()
	router := svc.appRoutes(r)

	svc.Settings = settings
	svc.Router = router
	svc.Handlers = Handlers{
		Projects: projects,
	}

	return svc, nil
}

// health is a simple health check handler that returns HTTP 200 OK.
func health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
