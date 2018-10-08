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
	return router
}
