package models

// Project is an object that contains files and data associated with a single project
type Project struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	PM       string `json:"pm"`
}

// AllProjects returns a list of all projects in the datastore
func (db *DB) AllProjects() ([]*Project, error) {

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
func (db *DB) CreateProject(p Project) (Project, error) {
	query := `INSERT INTO project (name, location, pm) VALUES ($1, $2, $3) RETURNING id, name, location, pm`

	err := db.QueryRowx(query, p.Name, p.Location, p.PM).StructScan(&p)
	return p, err
}
