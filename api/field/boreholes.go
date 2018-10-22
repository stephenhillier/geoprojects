package field

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

// BoreholeCreateRequest is the data a user should submit to create a borehole.
// A borehole can either be associated with an existing datapoint, or if a location
// is supplied, a datapoint will be created.
type BoreholeCreateRequest struct {
	Project   int64      `json:"project"`
	Program   NullInt64  `json:"program"`
	Datapoint NullInt64  `json:"datapoint"`
	Name      string     `json:"name"`
	StartDate NullDate   `json:"start_date" db:"start_date" schema:"start_date"`
	EndDate   NullDate   `json:"end_date" db:"end_date" schema:"end_date"`
	FieldEng  NullInt64  `json:"field_eng" db:"field_eng" schema:"field_eng"`
	Location  [2]float64 `json:"location"`
}

// BoreholeResponse is the data returned by the API after receiving a request for
// a borehole's details
// the FieldEng field is a string (users.username) instead of a primary key reference.
type BoreholeResponse struct {
	ID        int64      `json:"id"`
	Project   NullInt64  `json:"project"`
	Program   NullInt64  `json:"program"`
	Datapoint NullInt64  `json:"datapoint"`
	Name      string     `json:"name"`
	StartDate NullDate   `json:"start_date" db:"start_date"`
	EndDate   NullDate   `json:"end_date" db:"end_date"`
	FieldEng  NullString `json:"field_eng" db:"field_eng"`
}

// Borehole is drilled geotechnical test hole located at a Datapoint.
// There may be a number of samples/observations associated with one borehole.
type Borehole struct {
	ID        int64     `json:"id"`
	Project   NullInt64 `json:"project"`
	Program   NullInt64 `json:"program"`
	Datapoint int64     `json:"datapoint"`
	Name      string    `json:"name"`
	StartDate NullDate  `json:"start_date" db:"start_date"`
	EndDate   NullDate  `json:"end_date" db:"end_date"`
	FieldEng  NullInt64 `json:"field_eng" db:"field_eng"`
}

// boreholeOptions responds to OPTIONS requests (and pre-flight requests)
func (s *App) boreholeOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

// listBoreholes returns a list of boreholes for a specified project
// the project should be specified in a URL query string e.g. api/v1/boreholes?project=1
func (s *App) listBoreholes(w http.ResponseWriter, req *http.Request) {

	project := req.FormValue("project")

	var projectID int
	var err error

	// if a project was supplied in querystring, set projectID so that the db query can
	// list boreholes by project
	if project != "" {
		projectID, err = strconv.Atoi(project)
		if err != nil {
			// if project can't be converted to an int, make sure projectID is zero.
			// this ignores the ?project query if it's not a valid integer.
			projectID = 0
		}
	}

	boreholes, err := s.boreholes.ListBoreholes(projectID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	render.JSON(w, req, boreholes)
}

// createBorehole creates a new borehole
func (s *App) createBorehole(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	borehole := BoreholeCreateRequest{}
	err := decoder.Decode(&borehole)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	newBorehole, err := s.boreholes.CreateBorehole(borehole)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newBorehole)
}

// func getBorehole(w http.ResponseWriter, req *http.Request) {

// }
