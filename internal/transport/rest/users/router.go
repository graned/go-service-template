package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	endpoint "github.com/graned/go-service-template/internal/transport/rest/endpoint"
	"net/http"
)

func Routes(handler *Handler, validate *validator.Validate) chi.Router {
	r := chi.NewRouter()
	r.Get("/", endpoint.Serve(endpoint.Adapt(http.StatusOK, handler.Get)))
	r.Post("/", endpoint.Serve(
		endpoint.JSONBody(
			validate,
			http.StatusCreated,
			handler.Create,
		),
	),
	)
	r.Get("/{id}", endpoint.Serve(endpoint.Adapt(http.StatusOK, handler.Get)))

	return r
}
