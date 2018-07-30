package main

import (
	"github.com/go-chi/chi"
)

func (api *Server) routes() {
	// soil description parser endpoint (post a soil description as a
	// string, receive json with standardized soil properties)
	api.router.Post("/describe", api.Describe)

	// projects endpoints (list/create/retrieve/update/delete project records)
	api.router.Route("/projects", func(r chi.Router) {
		r.Get("/", api.ProjectsIndex)
	})
}
