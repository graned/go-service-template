package rest

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/graned/go-service-template/internal/transport/rest/users"
	domainuser "github.com/graned/go-service-template/internal/user"
	"log/slog"
	"net/http"
)

type Server struct {
	logger *slog.Logger
	server *http.Server
}

func New(
	address string,
	logger *slog.Logger,
	userService *domainuser.Service,
) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(requestLogger(logger))
	r.Use(middleware.Recoverer)

	validate := validator.New(
		validator.WithRequiredStructEnabled(),
	)
	usersHandlers := users.NewHandler(userService, validate)

	r.Get("/health", health)
	r.Mount("/users", users.Routes(usersHandlers))

	return &Server{
		logger: logger,
		server: &http.Server{
			Addr:    address,
			Handler: r,
		},
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Ok"))
}

func (s *Server) Run() error {
	s.logger.Info(
		"Rest API Server listening",
		"address", s.server.Addr,
	)

	err := s.server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down Rest API server")
	return s.server.Shutdown(ctx)
}
