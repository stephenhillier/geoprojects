# projects

This package and its sub-packages contain code for working with projects.

## http

The `http` package contains http handlers.  A service with access to handlers can be created with the http.NewProjectService function.

## repository

Package `repository` contains database operations.  There is a PostgresRepo type that implements all the methods in the earthworks.ProjectRepository interface (found in the repo's root).
