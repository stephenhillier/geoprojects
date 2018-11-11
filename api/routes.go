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
					r.Use(api.projectCtxMiddleware)
					r.Get("/", api.projectDetail)
					r.Options("/", api.singleProjectOptions)
					r.Delete("/", api.deleteProject)
				})
			})

			// Boreholes routes
			r.Route("/boreholes", func(r chi.Router) {
				r.Options("/", api.boreholeOptions)
				r.Get("/", api.listBoreholes)
				r.Post("/", api.createBorehole)
				r.Route("/{boreholeID}", func(r chi.Router) {
					r.Use(api.boreholeCtxMiddleware)
					r.Get("/", api.getBorehole)
					r.Get("/strata", api.listStrataByBorehole)
				})
			})

			// Soil strata routes
			r.Route("/strata", func(r chi.Router) {
				r.Options("/", api.strataOptions)
				r.Post("/", api.createStrata)
			})
		})

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(api.jwtAuthentication().Handler)
			// routes added here require authentication
		})
	})
	return r
}
