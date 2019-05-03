package http

import (
	"github.com/stephenhillier/geoprojects/earthworks"
	"github.com/stephenhillier/geoprojects/earthworks/db"
	projectsRepo "github.com/stephenhillier/geoprojects/earthworks/projects/repository"
)

// ProjectSvc is a service that provides methods for working with Projects
// http handlers will be available as methods e.g. projects.Create
type ProjectSvc struct {
	Repo     earthworks.ProjectRepository
	Settings earthworks.Settings
}

// NewProjectSvc returns a ProjectSvc with a handle for a database connection and app settings
func NewProjectSvc(db *db.Datastore, settings earthworks.Settings) *ProjectSvc {
	repo := projectsRepo.NewProjectsRepo(db)

	svc := &ProjectSvc{
		Repo:     repo,
		Settings: settings,
	}
	return svc
}
