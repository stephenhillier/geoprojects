package earthworks

import (
	"time"

	"github.com/stephenhillier/geoprojects/earthworks/db"
)

// FileRequest is a request that contains a file to be uploaded
type FileRequest struct {
	File      []byte
	Filename  string
	CreatedBy string
	Category  string
	Project   int
}

// File represents a file and some metadata about the file
type File struct {
	ID         int64       `json:"id"`
	Project    int64       `json:"project"`
	Filename   string      `json:"filename"`
	Category   string      `json:"category"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
	CreatedBy  string      `json:"created_by" db:"created_by"`
	Superseded bool        `json:"superseded"`
	Archived   db.NullDate `json:"archived,omitempty"`
}

// FileObject is a struct that contains a byte slice,
// for reading into and out of a database/file store.
type FileObject struct {
	File     []byte
	Filename string
}

// FileFilter allows searching for files based on different criteria (project, file category etc)
type FileFilter struct {
	ID       int64  `json:"id" schema:"id"`
	Project  int    `json:"project" schema:"project"`
	Category string `json:"category" schema:"category"`
	Archived bool   `json:"archived" schema:"archived"`
}

// FileRepository is a set of methods for interacting with files
type FileRepository interface {
	NewFile(FileRequest) (File, error)
	ListFiles(FileFilter) ([]File, error)
	GetFile(id int, project int) (FileObject, error)
	DeleteFile(id int, project int) error
	RestoreFile(id int, project int) error
}
