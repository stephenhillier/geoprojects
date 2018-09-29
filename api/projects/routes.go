package projects

import (
	"github.com/go-chi/chi"
)

func (s *App) routes() func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", s.listProjects)
		r.Options("/", s.projectOptions)
		r.Post("/", s.createProject)
		r.Route("/{projectID}", func(r chi.Router) {
			r.Get("/", s.projectDetail)
			r.Options("/", s.singleProjectOptions)
		})
	}
}
