package users

import (
	"github.com/go-chi/chi/v5"
)

func Routes(handler *Handler) chi.Router {
	r := chi.NewRouter()
	r.Get("/", handler.List)
	r.Post("/", handler.Create)
	r.Get("/{id}", handler.Get)

	return r
}
