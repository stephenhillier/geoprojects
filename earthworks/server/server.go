package server

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi"
	"github.com/stephenhillier/geoprojects/earthworks"
	boreholeHandlers "github.com/stephenhillier/geoprojects/earthworks/boreholes/http"
	"github.com/stephenhillier/geoprojects/earthworks/db"
	fileHandlers "github.com/stephenhillier/geoprojects/earthworks/files/http"
	projectsHandlers "github.com/stephenhillier/geoprojects/earthworks/projects/http"
)

// Service contains backing services and settings
// for the Earthworks backend app.
type Service struct {
	Config   *Config
	Settings earthworks.Settings
	Router   chi.Router
	Handlers Handlers
}

// Handlers contains http request handlers for each service
type Handlers struct {
	Projects  *projectsHandlers.ProjectSvc
	Boreholes *boreholeHandlers.BoreholeSvc
	Files     *fileHandlers.FileSvc
}

// Config holds server/database/auth service configuration
type Config struct {
	AuthCert         PEMCert // defined in auth.go
	AuthAudience     string
	AuthIssuer       string
	AuthJWKSEndpoint string
	DBConn           string
	DBDriver         string
	AuthGroupClaim   string
	AuthRoleClaim    string
}

// NewEarthworksService returns an EarthworksService object that ties
// together the various services that form the backend Earthworks application
func NewEarthworksService(datastore *db.Datastore, cnf *Config) (Service, error) {

	svc := Service{}

	settings := earthworks.Settings{
		DefaultPageLimit: 10,
		MaxPageLimit:     100,
	}

	projects := projectsHandlers.NewProjectSvc(datastore, settings)
	boreholes := boreholeHandlers.NewBoreholeSvc(datastore, settings)
	files := fileHandlers.NewFileSvc(datastore, settings)

	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	// register middleware
	r.Use(cors.Handler)
	r.Use(middleware.Logger)

	svc.Config = cnf
	svc.Settings = settings
	svc.Handlers = Handlers{
		Projects:  projects,
		Boreholes: boreholes,
		Files:     files,
	}

	router := svc.appRoutes(r)

	svc.Router = router

	return svc, nil
}
