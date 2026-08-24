package health

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/graned/go-service-template/internal/transport/rest/endpoint"
	"net/http"
)

func Routes(
	handler *Handler,
	validate *validator.Validate,
	errorHandler *endpoint.ErrorHandler,
) chi.Router {
	r := chi.NewRouter()
	r.Get("/", endpoint.Serve(errorHandler, endpoint.Adapt(http.StatusOK, handler.Get)))
	return r
}
