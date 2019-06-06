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
			r.Use(api.Config.JWTAuthentication().Handler)

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

					r.Get("/samples", api.Handlers.Boreholes.ListSamplesByProject)

					r.Route("/instrumentation", func(r chi.Router) {
						r.Get("/", api.Handlers.Instrumentation.List)
					})

					r.Route("/files", func(r chi.Router) {
						r.Post("/", api.Handlers.Files.NewFile)
						r.Get("/", api.Handlers.Files.ListFiles)
						r.Route("/{fileID}", func(r chi.Router) {
							r.Get("/", api.Handlers.Files.GetFile)
							r.Delete("/", api.Handlers.Files.DeleteFile)
						})
					})

					// Lab test routes

					r.Route("/lab/tests", func(r chi.Router) {
						r.Options("/", api.Handlers.Lab.LabTestOptions)
						r.Post("/", api.Handlers.Lab.CreateLabTest)
						r.Get("/", api.Handlers.Lab.ListLabTestsByProject)
						r.Route("/{labTestID}", func(r chi.Router) {
							r.Use(api.Handlers.Lab.LabTestCtxMiddleware)
							r.Options("/", api.Handlers.Lab.SingleLabTestOptions)
							r.Delete("/", api.Handlers.Lab.DeleteLabTest)
							r.Put("/", api.Handlers.Lab.PutLabTest)
							r.Get("/", api.Handlers.Lab.RetrieveLabTest)
							r.Route("/moisture", func(r chi.Router) {
								r.Get("/", api.Handlers.Lab.RetrieveMoistureTest)
								r.Post("/", api.Handlers.Lab.CreateMoistureTest)
								r.Put("/", api.Handlers.Lab.PutMoistureTest)
							})
							r.Route("/sieves", func(r chi.Router) {
								r.Post("/", api.Handlers.Lab.AddSieveToGSA)
								r.Route("/{sieveID}", func(r chi.Router) {
									r.Put("/", api.Handlers.Lab.PutSieve)
									r.Delete("/", api.Handlers.Lab.DeleteSieve)
								})
							})
						})
					})
				})
			})

			// Boreholes routes
			r.Route("/boreholes", func(r chi.Router) {
				r.Options("/", api.Handlers.Boreholes.Options)
				r.Get("/", api.Handlers.Boreholes.List)
				r.Post("/", api.Handlers.Boreholes.Create)
				r.Route("/{boreholeID}", func(r chi.Router) {
					r.Use(api.Handlers.Boreholes.BoreholeCtxMiddleware)
					r.Get("/", api.Handlers.Boreholes.Get)
					r.Get("/strata", api.Handlers.Boreholes.ListStrataByBorehole)
					r.Route("/samples", func(r chi.Router) {
						r.Options("/", api.Handlers.Boreholes.SampleOptions)
						r.Get("/", api.Handlers.Boreholes.ListSamplesByBorehole)
						r.Post("/", api.Handlers.Boreholes.CreateSample)
						r.Route("/{sampleID}", func(r chi.Router) {
							r.Use(api.Handlers.Boreholes.SampleCtxMiddleware)
							r.Get("/", api.Handlers.Boreholes.RetrieveSample)
							r.Put("/", api.Handlers.Boreholes.PutSample)
							r.Delete("/", api.Handlers.Boreholes.DeleteSample)
						})
					})

					r.Delete("/", api.Handlers.Boreholes.Delete)
				})
			})

			// Soil strata routes
			r.Route("/strata", func(r chi.Router) {
				r.Options("/", api.Handlers.Boreholes.StrataOptions)
				r.Post("/", api.Handlers.Boreholes.CreateStrata)
				r.Route("/{strataID}", func(r chi.Router) {
					r.Use(api.Handlers.Boreholes.StrataCtxMiddleware)
					r.Put("/", api.Handlers.Boreholes.PutStrata)
					r.Delete("/", api.Handlers.Boreholes.DeleteStrata)
				})
			})
		})
	})
	return r
}
