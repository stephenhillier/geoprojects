package main

import (
	"log"
	"strconv"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
)

// // ProjectsRepository is the set of methods available to a collection of projects
// type ProjectsRepository interface {
// 	AllProjects() ([]*Project, error)
// 	CreateProject(p Project) (Project, error)
// 	RetrieveProject(projectID int) (Project, error)
// 	DeleteProject(id int) error
// }

// AllProjects returns a list of all projects in the datastore
func (db *Datastore) AllProjects(name string, number string) ([]Project, error) {
	var err error

	query := `SELECT
			project.id,
			project.name,
			project.location,
			ST_AsBinary(project.default_coords) AS default_coords,
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
	}

	if len(number) > 0 {
		searchOnNum = true
		numParams++
		// check if name was also searched on
		if searchOnName {
			query = query + ` AND`
		} else {
			query = query + ` WHERE`
		}

		searchProjectNumber := ` CAST(project.id AS TEXT) LIKE $` + strconv.Itoa(numParams)

		query = query + searchProjectNumber
	}

	groupByQuery := ` GROUP BY project.id`
	query = query + groupByQuery

	if searchOnName && searchOnNum {
		err = db.Select(&projects, query, "%"+name+"%", number+"%")
	} else if searchOnName && !searchOnNum {
		err = db.Select(&projects, query, "%"+name+"%")
	} else if !searchOnName && searchOnNum {
		err = db.Select(&projects, query, number+"%")
	} else {
		err = db.Select(&projects, query)
	}

	if err != nil {
		log.Println(query, name, number)
		return []Project{}, err
	}
	return projects, nil
}

// CreateProject creates a new project record in the database
func (db *Datastore) CreateProject(p ProjectRequest) (Project, error) {
	query := `INSERT INTO project (name, number, client, pm, location, default_coords) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, number, client, pm, location`

	new := Project{}

	coords := orb.Point{p.DefaultCoords[0], p.DefaultCoords[1]}
	err := db.QueryRowx(query, p.Name, p.Number, p.Client, p.PM, p.Location, wkt.MarshalString(coords)).StructScan(&new)
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
							project.number,
							project.client,
							project.pm,
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
