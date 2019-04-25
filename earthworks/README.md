# Earthworks backend API

## Folders:

#### Project root files

The files in the project root (e.g. `projects.go`, `boreholes.go`) contain domain types that model the data managed by the Earthworks application.

#### Server package folders

The following special folders contain packages necessary for starting and operating the backend API server.

* `db/` is a package with a function `NewDB()` that returns a database connection handle, and also contains types for working with a Postgres database (e.g. `db.NullString`)
* `pkg/` has packages with utilities (e.g. the `pkg/gis` package contains types for working with spatial data)
* `cmd/` is where the server program (package main) lives. Other executable programs can be added here.

#### Domain/business folders

These folders contain business logic that the Earthworks backend API application uses to manage projects, geotechnical field data and laboratory data.

* `projects/`
* `boreholes/`
* `laboratory/`
