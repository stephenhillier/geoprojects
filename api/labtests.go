package main

import (
	"net/http"

	"github.com/go-chi/render"
)

// LabTest represents one lab test ordered on a given sample. A lab test should have
// at least one record of a specific test referencing it (e.g. a moisture content test record)
type LabTest struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"test_type" db:"type"`
	Sample int    `json:"sample" db:"sample"`
}

// listLabTests returns lab tests from a project.
// the project must be passed in the request context with the contextKey "projectCtx"
func (s *server) listLabTestsByBorehole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	project, ok := ctx.Value(projectCtx).(Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	labTests, err := s.datastore.ListLabTestsByProject(project.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, req, labTests)
}
