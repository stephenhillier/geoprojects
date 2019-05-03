package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stephenhillier/geoprojects/earthworks"
	"github.com/stephenhillier/soildesc"
)

// StrataOptions responds to preflight requests with allowed methods
func (svc *BoreholeSvc) StrataOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Allow", "GET, POST, PUT, DELETE, OPTIONS")
	return
}

// ListStrataByBorehole returns soil strata associated with a specified borehole.
// the borehole must be passed in the request context with the contextKey "boreholeCtx"
func (svc *BoreholeSvc) ListStrataByBorehole(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	borehole, ok := ctx.Value(earthworks.BoreholeCtx).(earthworks.BoreholeResponse)

	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	strata, err := svc.Repo.ListStrataByBorehole(borehole.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, req, strata)
}

// CreateStrata creates a soil strata/layer for a borehole
func (svc *BoreholeSvc) CreateStrata(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	strataReq := earthworks.StrataRequest{}
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

	strata := earthworks.Strata{
		Borehole:    strataReq.Borehole,
		Start:       strataReq.Start,
		End:         strataReq.End,
		Description: strataReq.Description,
		Soils:       strings.Join(parsedSoil, ", "),
		Moisture:    parsedDescription.Moisture,
		Consistency: parsedDescription.Consistency,
	}

	newStrata, err := svc.Repo.CreateStrata(strata)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusCreated)
	render.JSON(w, req, newStrata)
}

// PutStrata allows updating a strata record by making a PUT request to the strata's endpoint
func (svc *BoreholeSvc) PutStrata(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	strataReq := earthworks.StrataRequest{}
	err := decoder.Decode(&strataReq)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	ctx := req.Context()
	strata, ok := ctx.Value(earthworks.StrataCtx).(earthworks.Strata)
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

	updatedStrataData := earthworks.Strata{
		ID:          strata.ID,
		Borehole:    strataReq.Borehole,
		Start:       strataReq.Start,
		End:         strataReq.End,
		Description: strataReq.Description,
		Soils:       strings.Join(parsedSoil, ", "),
		Moisture:    parsedDescription.Moisture,
		Consistency: parsedDescription.Consistency,
	}

	updatedStrata, err := svc.Repo.UpdateStrata(updatedStrataData)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.Status(req, http.StatusOK)
	render.JSON(w, req, updatedStrata)
}

// DeleteStrata asks the datastore to delete a given strata record
func (svc *BoreholeSvc) DeleteStrata(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	strata, ok := ctx.Value(earthworks.StrataCtx).(earthworks.Strata)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	err := svc.Repo.DeleteStrata(strata.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.NoContent(w, req)
	return
}

// StrataCtxMiddleware is used by strata routes that have a strataID in the URL path.
// it finds the specified strata (returning 404 if the strata is not found) and adds it
// to the request context.
func (svc *BoreholeSvc) StrataCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strataID, err := strconv.Atoi(chi.URLParam(r, "strataID"))
		if err != nil {
			log.Println("strataID not supplied")
			http.Error(w, http.StatusText(404), 404)
			return
		}

		strata, err := svc.Repo.RetrieveStrata(strataID)
		if err != nil {
			log.Println("strata was not found in DB")
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), earthworks.StrataCtx, strata)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
