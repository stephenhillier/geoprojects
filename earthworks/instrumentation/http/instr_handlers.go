package http

import (
	"net/http"

	"github.com/go-chi/render"
)

// List displays all instrumentation for a given project
func (svc *InstrumentationSvc) List(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "hello")
}
