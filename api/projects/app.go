package projects

import "github.com/jmoiron/sqlx"

// NewApp takes a database and returns an app containing the
// the routes that this application has been configured to handle
func NewApp(db *sqlx.DB) *App {

	// create a Datastore with the supplied db and a Routes function
	app := &App{}
	app.repo = &Datastore{db}
	app.Routes = app.routes()

	return app
}
