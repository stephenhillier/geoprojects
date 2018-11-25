package main

import "log"

// // ProjectsRepository is the set of methods available to a collection of projects
// type ProjectsRepository interface {
// 	AllProjects() ([]*Project, error)
// 	CreateProject(p Project) (Project, error)
// 	RetrieveProject(projectID int) (Project, error)
// 	DeleteProject(id int) error
// }

// AllProjects returns a list of all projects in the datastore
func (db *Datastore) AllProjects(limit int, offset int) ([]Project, int, error) {
	countQuery := `SELECT count(id) FROM project`

	query := `SELECT
			project.id,
			project.name,
			project.location,
			COUNT(borehole.project) as borehole_count,
			ST_AsBinary(st_transform(st_centroid(st_union(st_transform(datapoint.location::geometry, 26910))), 4326)::geography) as centroid
		FROM project
		LEFT JOIN borehole ON (borehole.project = project.id)
		LEFT JOIN datapoint ON (borehole.datapoint = datapoint.id)
		GROUP BY project.id
		LIMIT $1 OFFSET $2
	`

	var count int
	projects := []Project{}

	err := db.Get(&count, countQuery)

	err = db.Select(&projects, query, limit, offset)
	if err != nil {
		return []Project{}, 0, err
	}

	return projects, count, nil
}

// CreateProject creates a new project record in the database
func (db *Datastore) CreateProject(p ProjectRequest) (Project, error) {
	query := `INSERT INTO project (name, location) VALUES ($1, $2) RETURNING id, name, location`

	new := Project{}

	log.Println("trying to insert")
	err := db.QueryRowx(query, p.Name, p.Location).StructScan(&new)
	if err != nil {
		return Project{}, err
	}
	return new, nil
}

// RetrieveProject fetches one project record from database (by project ID)
func (db *Datastore) RetrieveProject(projectID int) (Project, error) {
	p := Project{}
	query := `SELECT
							project.id,
							project.name,
							project.location,
							COUNT(borehole.project) as borehole_count,
							ST_AsBinary(st_transform(st_centroid(st_union(st_transform(datapoint.location::geometry, 26910))), 4326)::geography) as centroid
						FROM project
						LEFT JOIN borehole ON (borehole.project = project.id)
						LEFT JOIN datapoint ON (borehole.datapoint = datapoint.id)
						WHERE project.id=$1
						GROUP BY project.id
						`
	err := db.Get(&p, query, projectID)
	return p, err
}

// DeleteProject sets a project's expiry to the current time
func (db *Datastore) DeleteProject(id int) error {
	query := `DELETE FROM project WHERE id = $1`
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
