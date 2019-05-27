package http

import (
	"github.com/stephenhillier/geoprojects/earthworks"
	"github.com/stephenhillier/geoprojects/earthworks/db"
	labRepo "github.com/stephenhillier/geoprojects/earthworks/laboratory/repository"
)

// LabSvc is a service that provides methods for working with lab tests
// http handlers will be available as methods e.g. projects.Create
type LabSvc struct {
	Repo     earthworks.LabRepository
	Settings earthworks.Settings
}

// NewLabSvc returns a LabSvc with a handle for a database connection and app settings
func NewLabSvc(db *db.Datastore, settings earthworks.Settings) *LabSvc {
	repo := labRepo.NewLabRepo(db)

	svc := &LabSvc{
		Repo:     repo,
		Settings: settings,
	}
	return svc
}
