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

					r.Get("/samples", api.listSamplesByProject)

					r.Route("/files", func(r chi.Router) {
						r.Post("/", api.NewFile)
						r.Get("/", api.ListFiles)
						r.Get("/{fileID}", api.GetFile)
					})

					// Lab test routes

					r.Route("/lab/tests", func(r chi.Router) {
						r.Options("/", api.labTestOptions)
						r.Post("/", api.createLabTest)
						r.Get("/", api.listLabTestsByProject)
						r.Route("/{labTestID}", func(r chi.Router) {
							r.Use(api.labTestCtxMiddleware)
							r.Options("/", api.singleLabTestOptions)
							r.Delete("/", api.deleteLabTest)
							r.Put("/", api.putLabTest)
							r.Get("/", api.retrieveLabTest)
							r.Route("/moisture", func(r chi.Router) {
								r.Get("/", api.retrieveMoistureTest)
								r.Post("/", api.createMoistureTest)
								r.Put("/", api.putMoistureTest)
							})
							r.Route("/sieves", func(r chi.Router) {
								r.Post("/", api.addSieveToGSA)
								r.Route("/{sieveID}", func(r chi.Router) {
									r.Put("/", api.putSieve)
									r.Delete("/", api.deleteSieve)
								})
							})
						})
					})
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
