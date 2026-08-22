package rest

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	applogger "github.com/graned/go-service-template/internal/logger"
)

func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			requestID := middleware.GetReqID(r.Context())

			w.Header().Set(
				"X-Request-ID",
				requestID,
			)

			requestLogger := logger.With(
				"request_id", requestID,
			)

			ctx := applogger.WithContext(
				r.Context(),
				requestLogger,
			)

			requestLogger.InfoContext(
				ctx,
				"request started",
				"method", r.Method,
				"path", r.URL.Path,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}
