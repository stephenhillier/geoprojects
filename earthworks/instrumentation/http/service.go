package http

import (
	"github.com/stephenhillier/geoprojects/earthworks"
	"github.com/stephenhillier/geoprojects/earthworks/db"
	instrRepo "github.com/stephenhillier/geoprojects/earthworks/instrumentation/repository"
)

// InstrumentationSvc is a service that provides methods for working with instrumentation
// http handlers will be available as methods e.g. intrumentation.Create
type InstrumentationSvc struct {
	Repo     earthworks.InstrumentationRepository
	Settings earthworks.Settings
}

// NewInstrumentationSvc returns a InstrumentationSvc with a handle for a database connection and app settings
func NewInstrumentationSvc(db *db.Datastore, settings earthworks.Settings) *InstrumentationSvc {
	repo := instrRepo.NewInstrumentationRepo(db)

	svc := &InstrumentationSvc{
		Repo:     repo,
		Settings: settings,
	}
	return svc
}
