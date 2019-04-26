package repository

import (
	"fmt"
	"strings"

	"github.com/stephenhillier/geoprojects/earthworks"
	"github.com/stephenhillier/geoprojects/earthworks/db"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

// ProjectsRepository is the set of methods available for interacting with Projects records
type ProjectsRepository interface {
	AllProjects(name string, number string, search string) ([]earthworks.Project, error)
	CreateProject(p earthworks.ProjectRequest) (earthworks.Project, error)
	RetrieveProject(projectID int) (earthworks.Project, error)
	UpdateProject(id int, p earthworks.ProjectRequest) (earthworks.Project, error)
	DeleteProject(id int) error
}

// NewProjectsRepo returns a PostgresRepo with a database connection
// This method can be called with either a sqlx.DB or a sqlx.Tx (transaction)
func NewProjectsRepo(database *db.Datastore) *PostgresRepo {
	return &PostgresRepo{
		conn: database,
	}
}

// PostgresRepo has a database connection and methods to interact with projects in
// the database.
type PostgresRepo struct {
	conn *db.Datastore
}

// AllProjects returns a list of all projects in the datastore
func (repo *PostgresRepo) AllProjects(name string, number string, search string) ([]earthworks.Project, error) {
	projects := []earthworks.Project{}
	var err error

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q := psql.Select(`
			project.id,
			project.name,
			project.number,
			project.location,
			ST_AsBinary(project.default_coords) AS default_coords,
			COUNT(borehole.project) as borehole_count,
			ST_AsBinary(st_transform(st_centroid(st_union(st_transform(datapoint.location::geometry, 26910))), 4326)::geography) as centroid
			`).
		From("project").
		LeftJoin("borehole ON (borehole.project = project.id)").
		LeftJoin("datapoint ON (borehole.datapoint = datapoint.id)")

	if search != "" {
		q = q.Where(
			sq.Or{
				sq.Like{"LOWER(project.name)": strings.ToLower(fmt.Sprint("%", search, "%"))},
				sq.Like{"LOWER(project.number)": strings.ToLower(fmt.Sprint("%", search, "%"))},
			},
		)
	}

	q = q.GroupBy("project.id")

	projectsQuery, queryArgs, err := q.ToSql()

	if err != nil {
		return projects, err
	}

	err = repo.conn.Select(&projects, projectsQuery, queryArgs...)

	if err != nil {
		return []earthworks.Project{}, err
	}

	return projects, nil
}

// CreateProject creates a new project record in the database
func (repo *PostgresRepo) CreateProject(p earthworks.ProjectRequest) (earthworks.Project, error) {
	query := `INSERT INTO project (name, number, client, pm, location, default_coords) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, number, client, pm, location`

	new := earthworks.Project{}

	coords := orb.Point{p.DefaultCoords[0], p.DefaultCoords[1]}
	err := repo.conn.QueryRowx(query, p.Name, p.Number, p.Client, p.PM, p.Location, wkt.MarshalString(coords)).StructScan(&new)
	if err != nil {
		return earthworks.Project{}, err
	}
	return new, nil
}

// RetrieveProject fetches one project record from database (by project ID)
func (repo *PostgresRepo) RetrieveProject(projectID int) (earthworks.Project, error) {
	p := earthworks.Project{}
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
	err := repo.conn.Get(&p, query, projectID)
	return p, err
}

// DeleteProject sets a project's expiry to the current time
func (repo *PostgresRepo) DeleteProject(id int) error {
	query := `DELETE FROM project WHERE id = $1`
	_, err := repo.conn.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateProject updates the details of a project in the datastore
func (repo *PostgresRepo) UpdateProject(id int, p earthworks.ProjectRequest) (earthworks.Project, error) {
	query := `
	UPDATE project
	SET
		name = $1, number = $2, client = $3, pm = $4, location = $5
	WHERE id = $6
	RETURNING id, name, number, client, pm, location
	`

	proj := earthworks.Project{}

	err := repo.conn.QueryRowx(query, p.Name, p.Number, p.Client, p.PM, p.Location, id).StructScan(&proj)
	if err != nil {
		return earthworks.Project{}, err
	}
	return proj, nil
}
