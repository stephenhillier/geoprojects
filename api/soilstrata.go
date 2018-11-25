package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/stephenhillier/soildesc"
)

// Strata is a soil layer/stratum and contains information such as description and depth of the layer
type Strata struct {
	ID          int64   `json:"id"`
	Borehole    int64   `json:"borehole"`
	Start       float64 `json:"start" db:"start_depth"`
	End         float64 `json:"end" db:"end_depth"`
	Description string  `json:"description"`
	Soils       string  `json:"soils"`
	Moisture    string  `json:"moisture"`
	Consistency string  `json:"consistency"`
}

// StrataRequest is a struct containing fields required to create a new strata layer
type StrataRequest struct {
	ID          int64   `json:"id"`
	Borehole    int64   `json:"borehole,string"`
	Start       float64 `json:"start,string"`
	End         float64 `json:"end,string"`
	Description string  `json:"description"`
}

func (s *server) strataOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	return
}

// listStrataByBorehole returns soil strata associated with a specified borehole.
// the borehole must be passed in the request context with the contextKey "boreholeCtx"
func (s *server) listStrataByBorehole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	borehole, ok := ctx.Value(boreholeCtx).(BoreholeResponse)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	strata, err := s.datastore.ListStrataByBorehole(borehole.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, req, strata)
}

func (s *server) createStrata(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	strataReq := StrataRequest{}
	err := decoder.Decode(&strataReq)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	parsedSoil := soildesc.ParseSoilTerms(strataReq.Description)
	parsedDescription, err := soildesc.ParseDescription(strataReq.Description)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	strata := Strata{
		ID:          strataReq.ID,
		Borehole:    strataReq.Borehole,
		Start:       strataReq.Start,
		End:         strataReq.End,
		Description: strataReq.Description,
		Soils:       strings.Join(parsedSoil, ", "),
		Moisture:    parsedDescription.Moisture,
		Consistency: parsedDescription.Consistency,
	}

	newStrata, err := s.datastore.CreateStrata(strata)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newStrata)
}
