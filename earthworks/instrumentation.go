package earthworks

import (
	"time"

	"github.com/stephenhillier/geoprojects/earthworks/db"
	"github.com/stephenhillier/geoprojects/earthworks/pkg/gis"
)

// InstrumentationRepository has methods to create and store
// data about instrumentation in a database
type InstrumentationRepository interface {
	ListInstruments(project int) ([]*Instrument, error)
	CreateInstrument(instr InstrumentCreateRequest, project int) (Instrument, error)
	GetInstrument(id int64) (Instrument, error)
}

// Instrument is a data collection instrument that records
// various kinds of field data, such as temperature, settlement,
// inclination etc.
type Instrument struct {
	ID          int64             `json:"id"`
	Project     int64             `json:"project"`
	Program     db.NullInt64      `json:"program"`
	DeviceID    db.NullString     `json:"device_id" db:"device_id"`
	Name        string            `json:"name"`
	FieldEng    string            `json:"field_eng" db:"field_eng"`
	Type        db.NullString     `json:"type"` // the type of instrument. todo: code table
	Datapoint   int64             `json:"datapoint"`
	InstallDate time.Time         `json:"install_date" db:"install_date"`
	Location    gis.PointLocation `json:"location"`
}

// InstrumentCreateRequest is the data needed to create a new instrument
type InstrumentCreateRequest struct {
	Project     int64         `json:"project"`
	Name        string        `json:"name"`
	DeviceID    db.NullString `json:"device_id" db:"device_id"`
	Type        db.NullString `json:"type"`
	FieldEng    string        `json:"field_eng" db:"field_eng"`
	InstallDate db.NullDate   `json:"install_date" db:"install_date"`
	Datapoint   db.NullInt64  `json:"datapoint"`
	Location    [2]float64    `json:"location"`
}

// TimeSeriesData contains data in the form of a timestamp and a value.
type TimeSeriesData struct {
	ID         int64     `json:"id"`
	Instrument int64     `json:"instrument"`
	Timestamp  time.Time `json:"timestamp"`
	Value      float64   `json:"value"`
}

// InstrumentCtx is a context key for instrumentation
var InstrumentCtx struct{}
