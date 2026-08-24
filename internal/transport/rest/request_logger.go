package rest

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	applogger "github.com/graned/go-service-template/internal/logger"
)

func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			start := time.Now()

			requestID := middleware.GetReqID(r.Context())

			requestLogger := logger.With(
				"request_id", requestID,
			)

			ctx := applogger.WithContext(
				r.Context(),
				requestLogger,
			)

			w.Header().Set("X-Request-ID", requestID)

			ww := middleware.NewWrapResponseWriter(
				w,
				r.ProtoMajor,
			)
			requestLogger.DebugContext(
				ctx,
				"request started",
				"method", r.Method,
				"path", r.URL.Path,
			)

			next.ServeHTTP(
				ww,
				r.WithContext(ctx),
			)

			duration := time.Since(start)

			route := chi.RouteContext(
				r.Context(),
			).RoutePattern()

			attrs := []slog.Attr{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("route", route),
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.Float64("duration_ms", float64(duration)/float64(time.Millisecond)),
			}

			level := slog.LevelInfo

			switch {
			case ww.Status() >= 500:
				level = slog.LevelError

			case ww.Status() >= 400:
				level = slog.LevelWarn
			}

			logger.LogAttrs(
				ctx,
				level,
				"request completed",
				attrs...,
			)
		})
	}
}
