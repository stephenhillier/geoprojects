package projects

import (
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

// App represents an HTTP web application with a datastore, handlers and routes.
// Routes can be passed into a chi.Router Route() to provide an
// access point to the handlers in this app.
type App struct {
	repo   Repository
	Routes func(r chi.Router)
}

// NewApp takes a database and returns an App containing the
// routes that this application has been configured to handle
func NewApp(db *sqlx.DB) *App {

	// create a Datastore with the supplied db and a Routes function
	app := &App{}
	app.repo = &Datastore{db}
	app.Routes = app.makeRoutes()

	return app
}
