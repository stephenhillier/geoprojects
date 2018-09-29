package projects

import "github.com/jmoiron/sqlx"

// NewApp takes a database and returns an app containing the
// the routes that this application has been configured to handle
func NewApp(db *sqlx.DB) *App {

	// create a Datastore with the supplied db
	app := &App{}
	app.repo = &Datastore{db}
	app.Routes = app.routes()

	// create an
	return app
}
