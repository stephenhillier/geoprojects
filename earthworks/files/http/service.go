package http

import (
	"github.com/stephenhillier/geoprojects/earthworks"
	"github.com/stephenhillier/geoprojects/earthworks/db"
	fileRepo "github.com/stephenhillier/geoprojects/earthworks/files/repository"
)

// FileSvc is a service that provides methods for working with files
// http handlers will be available as methods e.g. files.Create
type FileSvc struct {
	Repo     earthworks.FileRepository
	Settings earthworks.Settings
}

// NewFileSvc returns a ProjectSvc with a handle for a database connection and app settings
func NewFileSvc(db *db.Datastore, settings earthworks.Settings) *FileSvc {
	repo := fileRepo.NewFileRepo(db)

	svc := &FileSvc{
		Repo:     repo,
		Settings: settings,
	}
	return svc
}
