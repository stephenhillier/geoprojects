package projects

import (
	"github.com/jmoiron/sqlx"
)

// Repository is the set of methods available to a collection of projects
type Repository interface {
	AllProjects() ([]*Project, error)
	CreateProject(p Project) (Project, error)
	RetrieveProject(projectID int) (Project, error)
	DeleteProject(id int) error
}

// Datastore is a database, provided by the API service when this app was initialized
type Datastore struct {
	*sqlx.DB
}

// AllProjects returns a list of all projects in the datastore
func (db *Datastore) AllProjects() ([]*Project, error) {

	query := `SELECT project.id, project.name, project.location FROM project WHERE expired_at IS NULL`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]*Project, 0)
	for rows.Next() {
		project := new(Project)
		err := rows.Scan(&project.ID, &project.Name, &project.Location)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

// CreateProject creates a new project record in the database
func (db *Datastore) CreateProject(p Project) (Project, error) {
	query := `INSERT INTO project (name, location) VALUES ($1, $2) RETURNING id, name, location`

	err := db.QueryRowx(query, p.Name, p.Location).StructScan(&p)
	return p, err
}

// RetrieveProject fetches one project record from database (by project ID)
func (db *Datastore) RetrieveProject(projectID int) (Project, error) {
	p := Project{}
	query := `SELECT project.id, project.name, project.location, COUNT(borehole.project) as borehole_count
						FROM project
						LEFT JOIN borehole ON (borehole.project = project.id)
						WHERE project.id=$1 AND project.expired_at IS NULL
						GROUP BY project.id
						`
	err := db.Get(&p, query, projectID)
	return p, err
}

// DeleteProject sets a project's expiry to the current time
func (db *Datastore) DeleteProject(id int) error {
	query := `UPDATE project SET expired_at = NOW() WHERE id = $1`
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
