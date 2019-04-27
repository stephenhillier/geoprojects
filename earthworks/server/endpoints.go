package server

import (
	"github.com/go-chi/chi"
)

func (api *Service) appRoutes(r chi.Router) chi.Router {

	// v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			// server health check
			r.Get("/health", health)
		})

		// 	// Authenticated routes
		r.Group(func(r chi.Router) {
			// r.Use(api.jwtAuthentication().Handler)

			// Projects routes
			r.Route("/projects", func(r chi.Router) {

				r.Get("/", api.Handlers.Projects.List)
				r.Options("/", api.Handlers.Projects.Options)
				r.Post("/", api.Handlers.Projects.Create)
				r.Route("/{projectID}", func(r chi.Router) {
					r.Use(api.Handlers.Projects.ProjectCtxMiddleware)
					r.Get("/", api.Handlers.Projects.Retrieve)
					r.Options("/", api.Handlers.Projects.ProjectDetailOptions)
					r.Delete("/", api.Handlers.Projects.Delete)
					r.Put("/", api.Handlers.Projects.Update)

					// 				r.Get("/samples", api.listSamplesByProject)

					// 				r.Route("/files", func(r chi.Router) {
					// 					r.Post("/", api.NewFile)
					// 					r.Get("/", api.ListFiles)
					// 					r.Route("/{fileID}", func(r chi.Router) {
					// 						r.Get("/", api.GetFile)
					// 						r.Delete("/", api.DeleteFile)
					// 					})
					// 				})

					// 				// Lab test routes

					// 				r.Route("/lab/tests", func(r chi.Router) {
					// 					r.Options("/", api.labTestOptions)
					// 					r.Post("/", api.createLabTest)
					// 					r.Get("/", api.listLabTestsByProject)
					// 					r.Route("/{labTestID}", func(r chi.Router) {
					// 						r.Use(api.labTestCtxMiddleware)
					// 						r.Options("/", api.singleLabTestOptions)
					// 						r.Delete("/", api.deleteLabTest)
					// 						r.Put("/", api.putLabTest)
					// 						r.Get("/", api.retrieveLabTest)
					// 						r.Route("/moisture", func(r chi.Router) {
					// 							r.Get("/", api.retrieveMoistureTest)
					// 							r.Post("/", api.createMoistureTest)
					// 							r.Put("/", api.putMoistureTest)
					// 						})
					// 						r.Route("/sieves", func(r chi.Router) {
					// 							r.Post("/", api.addSieveToGSA)
					// 							r.Route("/{sieveID}", func(r chi.Router) {
					// 								r.Put("/", api.putSieve)
					// 								r.Delete("/", api.deleteSieve)
					// 							})
					// 						})
					// 					})
					// 				})
				})
			})

			// 		// Boreholes routes
			// 		r.Route("/boreholes", func(r chi.Router) {
			// 			r.Options("/", boreholes.Options)
			// 			r.Get("/", boreholes.List)
			// 			r.Post("/", boreholes.Create)
			// 			r.Route("/{boreholeID}", func(r chi.Router) {
			// 				r.Use(boreholes.BoreholeCtxMiddleware)
			// 				r.Get("/", boreholes.Get)
			// 				r.Get("/strata", boreholes.ListStrataByBorehole)
			// 				r.Route("/samples", func(r chi.Router) {
			// 					r.Options("/", api.sampleOptions)
			// 					r.Get("/", api.listSamplesByBorehole)
			// 					r.Post("/", api.createSample)
			// 					r.Route("/{sampleID}", func(r chi.Router) {
			// 						r.Use(api.sampleCtxMiddleware)
			// 						r.Get("/", api.retrieveSample)
			// 						r.Put("/", api.putSample)
			// 						r.Delete("/", api.deleteSample)
			// 					})
			// 				})

			// 				r.Delete("/", boreholes.Delete)
			// 			})
			// 		})

			// 		// Soil strata routes
			// 		r.Route("/strata", func(r chi.Router) {
			// 			r.Options("/", boreholes.StrataOptions)
			// 			r.Post("/", boreholes.CreateStrata)
			// 			r.Route("/{strataID}", func(r chi.Router) {
			// 				r.Use(boreholes.StrataCtxMiddleware)
			// 				r.Put("/", boreholes.PutStrata)
			// 				r.Delete("/", boreholes.DeleteStrata)
			// 			})
			// 		})
		})
	})
	return r
}
