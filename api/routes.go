package main

import (
	"github.com/go-chi/chi"
)

func (api *server) routes() {

	// Public routes
	api.router.Group(func(r chi.Router) {
		// server health check
		r.Get("/health", api.health)
		r.Mount("/api/v1/projects", api.apps.projects.Routes)
		r.Mount("/api/v1/field", api.apps.field.Routes)
	})

	// Protected routes (requires JWT in header `Authorization: Bearer ___`)
	api.router.Group(func(r chi.Router) {
		r.Use(api.jwtAuthentication().Handler)
		// projects endpoints (list/create/retrieve/update/delete project records)
	})

}
