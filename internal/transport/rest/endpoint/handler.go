// Allows to define a strongly typed handler
package endpoint

import (
	applogger "github.com/graned/go-service-template/internal/logger"
	"log/slog"
	"net/http"
)

type Result struct {
	Status int
	Data   any
}

type Handler func(
	r *http.Request,
) (Result, error)

type TypedHandler[T any] func(
	r *http.Request,
) (T, error)

type StreamHandler func(
	r *http.Request,
) (StreamResult, error)

type BodyHandler[Request any, Response any] func(
	r *http.Request,
	body Request,
) (Response, error)

func Adapt[T any](
	status int,
	handler TypedHandler[T],
) Handler {
	return func(r *http.Request) (Result, error) {
		data, err := handler(r)
		if err != nil {
			return Result{}, err
		}

		return Result{
			Status: status,
			Data:   data,
		}, nil
	}
}

// Centrilized way to handle streams
func ServeStream(
	logger *slog.Logger,
	handler StreamHandler,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		result, err := handler(r)

		if err != nil {
			_ = writeError(w, err)
			return
		}

		w.Header().Set(
			"Content-Type",
			result.ContentType,
		)

		w.WriteHeader(result.Status)

		if err := result.Stream(
			r.Context(),
			w,
		); err != nil {
			logger.ErrorContext(
				r.Context(),
				"stream failed",
				"error", err,
			)
		}
	}
}

// Central HTTP handler adapter
func Serve(
	handler Handler,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		result, err := handler(r)

		logger := applogger.FromContext(r.Context())

		if err != nil {
			logger.ErrorContext(
				r.Context(),
				"request handler failed",
				"error", err,
			)

			if writeErr := writeError(w, err); writeErr != nil {
				logger.ErrorContext(
					r.Context(),
					"failed to write error response",
					"error", writeErr,
				)
			}

			return
		}

		if err := writeSuccess(w, result); err != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to write response",
				"error", err,
			)
		}
	}
}
