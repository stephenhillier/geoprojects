package projects

import (
	"github.com/go-chi/chi"
)

func (s *App) makeRoutes() chi.Router {
	router := chi.NewRouter()

	router.Get("/", s.listProjects)
	router.Options("/", s.projectOptions)
	router.Post("/", s.createProject)
	router.Route("/{projectID}", func(r chi.Router) {
		r.Get("/", s.projectDetail)
		r.Options("/", s.singleProjectOptions)
		r.Delete("/", s.deleteProject)
	})

	return router
}
