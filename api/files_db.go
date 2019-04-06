package main

import (
	sq "gopkg.in/Masterminds/squirrel.v1"
)

// NewFile adds a new file to the database and returns a file object
// containing info about the new record
func (db *Datastore) NewFile(file FileRequest) (File, error) {
	query := `INSERT INTO project_file (file, filename, project, category, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id, filename, project, created_by, created_at`

	new := File{}

	err := db.QueryRowx(query, file.File, file.Filename, file.Project, file.Category, file.CreatedBy).StructScan(&new)
	return new, err
}

// ListFiles take a FileFilter and returns files that match the filter
func (db *Datastore) ListFiles(filter FileFilter) ([]File, error) {

	var files []File

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

	err = db.Select(&files, fileQuery, queryArgs...)

	return files, err

}

// GetFile retrieves a project file by its ID
func (db *Datastore) GetFile(id int, project int) (FileObject, error) {
	file := FileObject{}
	q := `
		SELECT file, filename FROM project_file WHERE id=$1 AND project=$2 AND expired_at IS NULL
	`
	err := db.QueryRowx(q, id, project).StructScan(&file)
	return file, err
}

// DeleteFile sets the expired_at column to the current time, archiving the file.
func (db *Datastore) DeleteFile(id int, project int) error {
	f := File{}

	// First get the filename for the file to be archived, so we can archive all
	// versions of it
	nameQuery := `
		SELECT filename FROM project_file WHERE id=$1 AND project=$2 AND expired_at IS NULL
	`

	err := db.QueryRowx(nameQuery, id, project).StructScan(&f)
	if err != nil {
		return err
	}

	// now set all version of the file (for this project) to archived.
	q := `
		UPDATE project_file SET expired_at=NOW() WHERE filename=$1 AND project=$2
	`
	_, err = db.Exec(q, f.Filename, project)
	return err
}

// RestoreFile sets the expired_at column to null.
func (db *Datastore) RestoreFile(id int, project int) error {
	f := File{}

	// Get the filename of the file to be restored
	nameQuery := `
		SELECT filename FROM project_file WHERE id=$1 AND project=$2 AND expired_at IS NOT NULL
	`

	err := db.QueryRowx(nameQuery, id, project).StructScan(&f)
	if err != nil {
		return err
	}

	// un-expire all versions of the file by filename
	q := `
		UPDATE project_file SET expired_at=NULL WHERE filename=$1 AND project=$2
	`
	_, err = db.Exec(q, f.Filename, project)
	return err
}
