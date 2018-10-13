package projects

import (
	"github.com/jmoiron/sqlx"
)

// Project is an object that contains files and data associated with a single project
type Project struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Location      string `json:"location"`
	PM            string `json:"pm"`
	BoreholeCount int    `json:"borehole_count" db:"borehole_count"`
}

// Repository is the set of methods available to a collection of projects
type Repository interface {
	AllProjects() ([]*Project, error)
	CreateProject(p Project) (Project, error)
	RetrieveProject(projectID int) (Project, error)
}

// Datastore is a database, provided by the API service when this app was initialized
type Datastore struct {
	*sqlx.DB
}

// AllProjects returns a list of all projects in the datastore
func (db *Datastore) AllProjects() ([]*Project, error) {

	query := `SELECT project.id, project.name, project.location, users.username FROM project LEFT JOIN users ON project.pm=users.id;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]*Project, 0)
	for rows.Next() {
		project := new(Project)
		err := rows.Scan(&project.ID, &project.Name, &project.Location, &project.PM)
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
	query := `INSERT INTO project (name, location, pm) VALUES ($1, $2, $3) RETURNING id, name, location, pm`

	err := db.QueryRowx(query, p.Name, p.Location, p.PM).StructScan(&p)
	return p, err
}

// RetrieveProject fetches one project record from database (by project ID)
func (db *Datastore) RetrieveProject(projectID int) (Project, error) {
	p := Project{}
	query := `SELECT project.id, project.name, project.location, users.username AS pm, COUNT(borehole.project) as borehole_count
						FROM project
						LEFT JOIN users ON project.pm=users.id 
						LEFT JOIN borehole ON (borehole.project = project.id)
						WHERE project.id=$1
						GROUP BY project.id, users.username
						`
	err := db.Get(&p, query, projectID)
	return p, err
}
