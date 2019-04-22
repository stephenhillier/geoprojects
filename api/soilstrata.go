package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	boreholev1 "github.com/stephenhillier/geoprojects/api/boreholes/model"
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
	Borehole    int64   `json:"borehole,string"`
	Start       float64 `json:"start,string"`
	End         float64 `json:"end,string"`
	Description string  `json:"description"`
}

func (s *server) strataOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, PUT, DELETE, OPTIONS")
	return
}

// listStrataByBorehole returns soil strata associated with a specified borehole.
// the borehole must be passed in the request context with the contextKey "boreholeCtx"
func (s *server) listStrataByBorehole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	borehole, ok := ctx.Value(boreholev1.BoreholeCtx).(boreholev1.BoreholeResponse)

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

// putStrata allows updating a strata record by making a PUT request to the strata's endpoint
func (s *server) putStrata(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	strataReq := StrataRequest{}
	err := decoder.Decode(&strataReq)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	ctx := req.Context()
	strata, ok := ctx.Value(strataCtx).(Strata)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	parsedSoil := soildesc.ParseSoilTerms(strataReq.Description)
	parsedDescription, err := soildesc.ParseDescription(strataReq.Description)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	updatedStrataData := Strata{
		ID:          strata.ID,
		Borehole:    strataReq.Borehole,
		Start:       strataReq.Start,
		End:         strataReq.End,
		Description: strataReq.Description,
		Soils:       strings.Join(parsedSoil, ", "),
		Moisture:    parsedDescription.Moisture,
		Consistency: parsedDescription.Consistency,
	}

	updatedStrata, err := s.datastore.UpdateStrata(updatedStrataData)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusOK)
	render.JSON(w, req, updatedStrata)
}

// deleteStrata asks the datastore to delete a given strata record
func (s *server) deleteStrata(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	strata, ok := ctx.Value(strataCtx).(Strata)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	err := s.datastore.DeleteStrata(strata.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.NoContent(w, req)
	return
}

// strataCtxMiddleware is used by strata routes that have a strataID in the URL path.
// it finds the specified strata (returning 404 if the strata is not found) and adds it
// to the request context.
func (s *server) strataCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strataID, err := strconv.Atoi(chi.URLParam(r, "strataID"))
		if err != nil {
			log.Println("strataID not supplied")
			http.Error(w, http.StatusText(404), 404)
			return
		}

		strata, err := s.datastore.RetrieveStrata(strataID)
		if err != nil {
			log.Println("strata was not found in DB")
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), strataCtx, strata)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
