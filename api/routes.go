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
		})

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(api.jwtAuthentication().Handler)

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
					r.Route("/samples", func(r chi.Router) {
						r.Options("/", api.sampleOptions)
						r.Get("/", api.listSamplesByBorehole)
						r.Post("/", api.createSample)
						r.Route("/{sampleID}", func(r chi.Router) {
							r.Use(api.sampleCtxMiddleware)
							r.Get("/", api.retrieveSample)
							r.Put("/", api.putSample)
							r.Delete("/", api.deleteSample)
						})
					})

					r.Delete("/", api.deleteBorehole)
				})
			})

			// Soil strata routes
			r.Route("/strata", func(r chi.Router) {
				r.Options("/", api.strataOptions)
				r.Post("/", api.createStrata)
				r.Route("/{strataID}", func(r chi.Router) {
					r.Use(api.strataCtxMiddleware)
					r.Put("/", api.putStrata)
					r.Delete("/", api.deleteStrata)
				})
			})
		})
	})
	return r
}
