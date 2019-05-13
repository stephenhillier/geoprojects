package repository

import (
	sq "gopkg.in/Masterminds/squirrel.v1"

	"github.com/stephenhillier/geoprojects/earthworks"
	"github.com/stephenhillier/geoprojects/earthworks/db"
)

// NewFileRepo returns a PostgresRepo with a database connection
func NewFileRepo(database *db.Datastore) *PostgresRepo {
	return &PostgresRepo{
		database,
	}
}

// PostgresRepo has a database connection and methods to interact with files in
// a Postgres database.
type PostgresRepo struct {
	db *db.Datastore // a database where file metadata is stored.
}

// NewFile adds a new file to the database and returns a file object
// containing info about the new record
func (repo *PostgresRepo) NewFile(file earthworks.FileRequest) (earthworks.File, error) {
	query := `INSERT INTO project_file (file, filename, project, category, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id, filename, project, created_by, created_at`

	new := earthworks.File{}

	err := repo.db.QueryRowx(query, file.File, file.Filename, file.Project, file.Category, file.CreatedBy).StructScan(&new)
	return new, err
}

// ListFiles take a FileFilter and returns files that match the filter
func (repo *PostgresRepo) ListFiles(filter earthworks.FileFilter) ([]earthworks.File, error) {

	var files []earthworks.File

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q := psql.Select(`
			id, filename, created_by, created_at, category	
		`).
		From("project_file")

	if filter.Category != "" {
		q = q.Where(sq.Eq{"category": filter.Category})
	}

	if filter.Project != 0 {
		q = q.Where(sq.Eq{"project": filter.Project})
	}

	// filter out archived/deleted files.  This is the default.
	if !filter.Archived {
		q = q.Where(sq.Eq{"expired_at": nil})
	}

	q = q.OrderBy("filename ASC, created_at DESC")

	fileQuery, queryArgs, err := q.ToSql()

	if err != nil {
		return files, err
	}

	err = repo.db.Select(&files, fileQuery, queryArgs...)

	return files, err

}

// GetFile retrieves a project file by its ID
func (repo *PostgresRepo) GetFile(id int, project int) (earthworks.FileObject, error) {
	file := earthworks.FileObject{}
	q := `
		SELECT file, filename FROM project_file WHERE id=$1 AND project=$2 AND expired_at IS NULL
	`
	err := repo.db.QueryRowx(q, id, project).StructScan(&file)
	return file, err
}

// DeleteFile sets the expired_at column to the current time, archiving the file.
func (repo *PostgresRepo) DeleteFile(id int, project int) error {
	f := earthworks.File{}

	// First get the filename for the file to be archived, so we can archive all
	// versions of it
	nameQuery := `
		SELECT filename FROM project_file WHERE id=$1 AND project=$2 AND expired_at IS NULL
	`

	err := repo.db.QueryRowx(nameQuery, id, project).StructScan(&f)
	if err != nil {
		return err
	}

	// now set all version of the file (for this project) to archived.
	q := `
		UPDATE project_file SET expired_at=NOW() WHERE filename=$1 AND project=$2
	`
	_, err = repo.db.Exec(q, f.Filename, project)
	return err
}

// RestoreFile sets the expired_at column to null.
func (repo *PostgresRepo) RestoreFile(id int, project int) error {
	f := earthworks.File{}

	// Get the filename of the file to be restored
	nameQuery := `
		SELECT filename FROM project_file WHERE id=$1 AND project=$2 AND expired_at IS NOT NULL
	`

	err := repo.db.QueryRowx(nameQuery, id, project).StructScan(&f)
	if err != nil {
		return err
	}

	// un-expire all versions of the file by filename
	q := `
		UPDATE project_file SET expired_at=NULL WHERE filename=$1 AND project=$2
	`
	_, err = repo.db.Exec(q, f.Filename, project)
	return err
}
