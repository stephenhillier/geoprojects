package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stephenhillier/geoprojects/earthworks"
)

// List displays all instrumentation for a given project
func (svc *InstrumentationSvc) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project, ok := ctx.Value(earthworks.ContextKey{Name: "ProjectContext"}).(earthworks.Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	instruments, err := svc.Repo.ListInstruments(project.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	render.JSON(w, r, instruments)

}

// Create handles a POST request and creates a new instrument in the repo
func (svc *InstrumentationSvc) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	instrRequest := earthworks.InstrumentCreateRequest{}
	err := decoder.Decode(&instrRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	// get project context
	ctx := r.Context()
	project, ok := ctx.Value(earthworks.ContextKey{Name: "ProjectContext"}).(earthworks.Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	newInstrument, err := svc.Repo.CreateInstrument(instrRequest, project.ID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, newInstrument)
}

// Get handles a GET request and responds with details about a single instrument
func (svc *InstrumentationSvc) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project, ok := ctx.Value(earthworks.ContextKey{Name: "ProjectContext"}).(earthworks.Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	instr, ok := ctx.Value(earthworks.InstrumentCtx).(earthworks.Instrument)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// return 404 if there's a mismatch between the project in the URL and the instrument's project.
	if int64(project.ID) != instr.Project {
		http.Error(w, http.StatusText(404), 404)
	}

	render.JSON(w, r, instr)

}

// InstrumentCtxMiddleware is used by instrument routes that have an instrument ID in the url path.
func (svc *InstrumentationSvc) InstrumentCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		instrumentID, err := strconv.ParseInt(chi.URLParam(r, "instrumentID"), 10, 64)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		instrument, err := svc.Repo.GetInstrument(instrumentID)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), earthworks.InstrumentCtx, instrument)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// PostTimeSeriesData handles an incoming request that contains a value and a timestamp.
// These requests should normally originate from an instrumentation device as the device
// records and reports readings.
func (svc *InstrumentationSvc) PostTimeSeriesData(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	dataRequest := earthworks.TimeSeriesData{}
	err := decoder.Decode(&dataRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	// create a "zero value" time object.  The request timestamp will be "zero"
	// if a timestamp was not provided along with the data.  If that is the case,
	// create a new timestamp now.  Ideally, the data should come with a timestamp
	// in case there were any delays in reporting the data, but in the worst case
	// scenario we still always need a time.
	zeroTime := time.Time{}
	if dataRequest.Timestamp == zeroTime {
		dataRequest.Timestamp = time.Now()
	}

	new, err := svc.Repo.PostTimeSeriesData(dataRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, new)
}

// GetTimeSeriesData handles a request for time series based instrumentation data
// (e.g., data with a value and a timestamp)
func (svc *InstrumentationSvc) GetTimeSeriesData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	project, ok := ctx.Value(earthworks.ContextKey{Name: "ProjectContext"}).(earthworks.Project)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	instr, ok := ctx.Value(earthworks.InstrumentCtx).(earthworks.Instrument)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	// return 404 if there's a mismatch between the project in the URL and the instrument's project.
	if int64(project.ID) != instr.Project {
		http.Error(w, http.StatusText(404), 404)
	}

	data := []*earthworks.TimeSeriesData{}

	deviceID, err := instr.DeviceID.Value()
	id, ok := deviceID.(string)
	if err != nil || !ok {
		// no device ID, return empty dataset now
		// Client already knows that there is no device_id so should be able to
		// warn user if required.
		render.JSON(w, r, data)
		return
	}

	data, err = svc.Repo.FetchTimeSeriesData(id)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	render.JSON(w, r, data)
}
