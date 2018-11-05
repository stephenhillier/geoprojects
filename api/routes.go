package main

import (
	"github.com/go-chi/chi"
)

func (api *server) appRoutes(r chi.Router) chi.Router {

	// v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			// server health check
			r.Get("/health", api.health)

			// Projects routes
			r.Route("/projects", func(r chi.Router) {
				r.Get("/", api.listProjects)
				r.Options("/", api.projectOptions)
				r.Post("/", api.createProject)
				r.Route("/{projectID}", func(r chi.Router) {
					r.Get("/", api.projectDetail)
					r.Options("/", api.singleProjectOptions)
					r.Delete("/", api.deleteProject)
				})
			})

			// Programs routes
			r.Route("/programs", func(r chi.Router) {
				r.Get("/", api.listPrograms)
				r.Options("/", api.programOptions)
				r.Post("/", api.createProgram)
				// r.Route("/{programID}", func(r chi.Router) {
				// 	r.Get("/", api.programDetail)
				// 	r.Options("/", api.singleProgramOptions)
				// })
			})

			// Boreholes routes
			r.Route("/boreholes", func(r chi.Router) {
				r.Options("/", api.boreholeOptions)
				r.Get("/", api.listBoreholes)
				r.Post("/", api.createBorehole)
			})
		})

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(api.jwtAuthentication().Handler)
			// projects endpoints (list/create/retrieve/update/delete project records)
		})
	})
	return r
}
