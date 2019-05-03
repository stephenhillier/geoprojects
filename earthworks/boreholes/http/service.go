package http

import (
	"github.com/stephenhillier/geoprojects/earthworks"
	boreholeRepo "github.com/stephenhillier/geoprojects/earthworks/boreholes/repository"
	"github.com/stephenhillier/geoprojects/earthworks/db"
)

// BoreholeSvc is a service that provides methods for working with Boreholes
// http handlers will be available as methods e.g. projects.Create
type BoreholeSvc struct {
	Repo     earthworks.BoreholeRepository
	Settings earthworks.Settings
}

// NewBoreholeSvc returns a ProjectSvc with a handle for a database connection and app settings
func NewBoreholeSvc(db *db.Datastore, settings earthworks.Settings) *BoreholeSvc {
	repo := boreholeRepo.NewBoreholeRepo(db)

	svc := &BoreholeSvc{
		Repo:     repo,
		Settings: settings,
	}
	return svc
}
