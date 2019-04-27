package earthworks

import (
	"github.com/paulmach/orb"
	"github.com/stephenhillier/geoprojects/earthworks/db"
	"github.com/stephenhillier/geoprojects/earthworks/pkg/gis"
)

// Project represents an engineering project. It holds files and data associated with a single project
type Project struct {
	ID int `json:"id"`

	// Number refers to a Project Number.  This is a typical industry term but it often contains letters
	// (to indicate the region, department etc)
	Number           string            `json:"number"`
	Name             string            `json:"name"`
	Location         string            `json:"location"`
	Client           string            `json:"client"`
	PM               string            `json:"pm"`
	BoreholeCount    int               `json:"borehole_count" db:"borehole_count"`
	CentroidLocation gis.PointLocation `json:"centroid" db:"centroid"`
	DefaultCoords    gis.PointLocation `json:"default_coords" db:"default_coords"`
}

// ProjectRequest is the set of data required to accept a request for a new project
type ProjectRequest struct {
	Name          string     `json:"name"`
	Number        string     `json:"number"`
	Client        string     `json:"client"`
	PM            string     `json:"pm"`
	Location      string     `json:"location"`
	DefaultCoords [2]float64 `json:"default_coords"`
}

// Datapoint is a geographic point that represents the location where data
// was gathered or where boreholes/instruments are located. A single datapoint
// may have a variety of data or records associated with it.
// Boreholes/instruments are organized this way (as children of a Datapoint) in
// order to reflect that drilling/subsurface sampling/instrumentation would
// generally have been performed at the same physical location.
type Datapoint struct {
	ID        db.NullInt64   `json:"id"`
	Location  orb.Point      `json:"location"`
	Elevation db.NullFloat64 `json:"elevation"`
}

// ProjectCtx is a context key for a project
var ProjectCtx struct{}

// ProjectRepository is the set of methods available for interacting with Projects records
type ProjectRepository interface {
	AllProjects(name string, number string, search string) ([]Project, error)
	CreateProject(p ProjectRequest) (Project, error)
	RetrieveProject(projectID int) (Project, error)
	UpdateProject(id int, p ProjectRequest) (Project, error)
	DeleteProject(id int) error
}
