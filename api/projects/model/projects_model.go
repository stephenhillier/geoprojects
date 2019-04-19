package model

import "github.com/stephenhillier/geoprojects/api/gis"

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
