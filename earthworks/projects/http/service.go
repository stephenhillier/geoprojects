package http

import (
	"github.com/stephenhillier/geoprojects/earthworks"
	"github.com/stephenhillier/geoprojects/earthworks/db"
	projectsRepo "github.com/stephenhillier/geoprojects/earthworks/projects/repository"
)

// NewProjectSvc returns a PRojectSvc with a handle for a database connection and app settings
func NewProjectSvc(db *db.Datastore, settings earthworks.Settings) *ProjectSvc {
	repo := projectsRepo.NewProjectsRepo(db)

	svc := &ProjectSvc{
		Repo:     repo,
		Settings: settings,
	}
	return svc
}
