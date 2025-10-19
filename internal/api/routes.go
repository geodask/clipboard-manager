package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/health", h.Health)

	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/history", func(r chi.Router) {
			r.Get("/", h.GetHistory)
			r.Get("/{id}", h.GetEntry)

			r.Delete("/", h.ClearHistory)
			r.Delete("/{id}", h.DeleteEntry)
		})

		r.Post("/entries", h.CreateEntry)

		r.Get("/search", h.Search)

		r.Get("/stats", h.GetStats)

	})

	return r

}
