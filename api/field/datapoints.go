package field

import (
	"github.com/paulmach/orb"
)

// Datapoint is a geographic point that represents the location where data
// was gathered or where boreholes/instruments are located. A single datapoint
// may have a variety of data or records associated with it.
// Boreholes/instruments are organized this way (as children of a Datapoint) in
// order to reflect that drilling/subsurface sampling/instrumentation would
// generally have been performed at the same physical location.
type Datapoint struct {
	ID       NullInt64 `json:"id"`
	Location orb.Point `json:"location"`
}
