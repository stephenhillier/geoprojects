package models

// Project is an object that contains files and data associated with a single project
type Project struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// AllProjects returns a list of all projects in the datastore
func (db *DB) AllProjects() ([]*Project, error) {

	query := `SELECT * FROM project`

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
