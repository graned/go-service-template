package endpoint

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type ErrorHandler struct {
	logger *slog.Logger
}

func NewErrorHandler(
	logger *slog.Logger,
) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

func (h *ErrorHandler) Handle(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	if writeErr := WriteError(w, err); writeErr != nil {
		h.logger.ErrorContext(
			r.Context(),
			"failed to write error response",
			"error", writeErr,
		)
	}
}

func recoverer(
	errorHandler *ErrorHandler,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			defer func() {
				if recovered := recover(); recovered != nil {
					// log panic error
					errorHandler.Handle(w, r, ApiInternalServerError(nil))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func RegisterErrorHandlers(
	router *chi.Mux,
	errorHandler *ErrorHandler,
) {
	router.Use(recoverer(errorHandler))
	router.NotFound(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		errorHandler.Handle(w, r, ApiRouteNotFound(nil))
	})

	router.MethodNotAllowed(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		errorHandler.Handle(w, r, ApiMethodNotAllowedError(nil))
	})
}
