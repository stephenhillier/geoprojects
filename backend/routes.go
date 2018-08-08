package main

import (
	"github.com/go-chi/chi"
)

func (api *Server) routes() {

	// Public routes
	api.router.Group(func(r chi.Router) {
		// server health check
		r.Get("/health", api.health)

		// soil description parser endpoint (post a soil description as a
		// string, receive json with standardized soil properties)
		r.Post("/describe", api.Describe)
	})

	// Protected routes (requires JWT in header `Authorization: Bearer ___`)
	api.router.Group(func(r chi.Router) {
		r.Use(jwtAuthentication().Handler)
		// projects endpoints (list/create/retrieve/update/delete project records)
		r.Route("/projects", func(r chi.Router) {
			r.Get("/", api.ProjectsIndex)
			r.Options("/", api.ProjectOpts)
			r.Post("/", api.ProjectPost)
		})
	})

}
