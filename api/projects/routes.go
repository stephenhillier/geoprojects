package projects

import (
	"github.com/go-chi/chi"
)

// NewRouter returns a router for projects endpoints
func NewRouter(projects ProjectSvc) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", projects.List)
	r.Options("/", projects.Options)
	r.Post("/", projects.Create)
	r.Route("/{projectID}", func(r chi.Router) {
		r.Use(projects.projectCtxMiddleware)
		r.Get("/", projects.Retrieve)
		r.Options("/", projects.ProjectDetailOptions)
		r.Delete("/", projects.Delete)
		r.Put("/", projects.Update)
	})
	return r
}
