package field

import (
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

// App represents an HTTP web application with a datastore, handlers and routes.
// Routes can be passed into a chi.Router Route() to provide an
// access point to the handlers in this app.
type App struct {
	programs  ProgramRepository
	boreholes BoreholeRepository
	Routes    chi.Router
}

// datastore is the database containing records related to the field module
// e.g. field programs, datapoints, boreholes
type datastore struct {
	*sqlx.DB
}

// NewApp takes a database and returns an App containing the
// routes that this application has been configured to handle
func NewApp(db *sqlx.DB) *App {
	app := &App{}
	app.programs = &datastore{db}
	app.boreholes = &datastore{db}
	app.Routes = app.makeRoutes()
	return app
}
