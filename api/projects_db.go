package main

import (
	"log"
	"strconv"
)

// // ProjectsRepository is the set of methods available to a collection of projects
// type ProjectsRepository interface {
// 	AllProjects() ([]*Project, error)
// 	CreateProject(p Project) (Project, error)
// 	RetrieveProject(projectID int) (Project, error)
// 	DeleteProject(id int) error
// }

// AllProjects returns a list of all projects in the datastore
func (db *Datastore) AllProjects(limit int, offset int, name string, number string) ([]Project, int, error) {
	var count int
	var err error

	// count and query errors
	var errC error
	var errQ error

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
	`

	projects := []Project{}

	numParams := 0

	searchOnName := false
	searchOnNum := false

	// append WHERE statements
	// TODO: figure out a better way to do this
	// note: all input is parameterized.
	if len(name) > 0 {
		numParams++
		searchOnName = true

		searchProjectName := ` WHERE project.name ILIKE $` + strconv.Itoa(numParams)

		query = query + searchProjectName
		countQuery = countQuery + searchProjectName
	}

	if len(number) > 0 {
		searchOnNum = true
		numParams++
		// check if name was also searched on
		if searchOnName {
			query = query + ` AND`
			countQuery = countQuery + ` AND`
		} else {
			query = query + ` WHERE`
			countQuery = countQuery + ` WHERE`
		}

		searchProjectNumber := ` CAST(project.id AS TEXT) LIKE $` + strconv.Itoa(numParams)

		query = query + searchProjectNumber
		countQuery = countQuery + searchProjectNumber
	}

	limitQuery := ` GROUP BY project.id LIMIT $` + strconv.Itoa(numParams+1) + ` OFFSET $` + strconv.Itoa(numParams+2) + ` `
	query = query + limitQuery

	if searchOnName && searchOnNum {
		errC = db.Get(&count, countQuery, "%"+name+"%", number+"%")
		errQ = db.Select(&projects, query, "%"+name+"%", number+"%", limit, offset)
	} else if searchOnName && !searchOnNum {
		errC = db.Get(&count, countQuery, "%"+name+"%")
		errQ = db.Select(&projects, query, "%"+name+"%", limit, offset)
	} else if !searchOnName && searchOnNum {
		errC = db.Get(&count, countQuery, number+"%")
		errQ = db.Select(&projects, query, number+"%", limit, offset)
	} else {
		errC = db.Get(&count, countQuery)
		errQ = db.Select(&projects, query, limit, offset)
	}

	if errQ != nil {
		log.Println(query, name, number, limit, offset)
		return []Project{}, 0, err
	}
	if errC != nil {
		log.Println("error counting results:", err)
	}

	return projects, count, nil
}

// CreateProject creates a new project record in the database
func (db *Datastore) CreateProject(p ProjectRequest) (Project, error) {
	query := `INSERT INTO project (name, location) VALUES ($1, $2) RETURNING id, name, location`

	new := Project{}

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
