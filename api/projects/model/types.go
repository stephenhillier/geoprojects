package model

import (
	"encoding/json"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
)

// ProjectCtx is a context key for a project
var ProjectCtx struct{}

// PointLocation is an alias for orb.Point
type PointLocation struct {
	orb.Point
}

// MarshalJSON represents PointLocation (orb.Point) as JSON
func (v PointLocation) MarshalJSON() ([]byte, error) {
	return json.Marshal([]float64{v.Lat(), v.Lon()})
}

// Scan allows scanning of PostGIS binary locations
func (v *PointLocation) Scan(src interface{}) error {

	if src == nil {
		emptyPoint := orb.Point{}
		*v = PointLocation{emptyPoint}
		return nil
	}

	var err error
	source := src.([]byte)

	geom, err := wkb.Unmarshal(source)

	if err != nil {
		return err
	}

	point := geom.(orb.Point)

	*v = PointLocation{point}

	return nil
}
