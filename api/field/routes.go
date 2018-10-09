package field

import (
	"github.com/go-chi/chi"
)

func (s *App) makeRoutes() chi.Router {
	router := chi.NewRouter()

	router.Route("/programs", func(r chi.Router) {
		r.Get("/", s.listPrograms)
		r.Options("/", s.programOptions)
		r.Post("/", s.createProgram)
		// r.Route("/{programID}", func(r chi.Router) {
		// 	r.Get("/", s.programDetail)
		// 	r.Options("/", s.singleProgramOptions)
		// })
	})
	router.Route("/boreholes", func(r chi.Router) {
		r.Options("/", s.boreholeOptions)
		r.Get("/", s.listBoreholes)
		r.Post("/", s.createBorehole)
	})
	return router
}
