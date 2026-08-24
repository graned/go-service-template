package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

	"github.com/graned/go-service-template/internal/app"
	"github.com/graned/go-service-template/internal/transport"
	"github.com/graned/go-service-template/internal/transport/rest/endpoint"
	"github.com/graned/go-service-template/internal/transport/rest/health"
	"github.com/graned/go-service-template/internal/transport/rest/users"
	domainuser "github.com/graned/go-service-template/internal/user"
)

type Server struct {
	transport.Runtime
	server *http.Server
}

func New(
	runtime transport.Runtime,
	userService *domainuser.Service,
) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(requestLogger(runtime.Logger))

	validate := validator.New(
		validator.WithRequiredStructEnabled(),
	)
	usersHandlers := users.NewHandler(userService)
	healthHandlers := health.NewHandler()

	errorHandler := endpoint.NewErrorHandler(runtime.Logger)

	endpoint.RegisterErrorHandlers(r, errorHandler)

	r.Mount("/health", health.Routes(healthHandlers, validate, errorHandler))
	r.Mount("/users", users.Routes(usersHandlers, validate, errorHandler))

	return &Server{
		Runtime: runtime,

		server: &http.Server{
			Addr:    runtime.Address,
			Handler: r,
		},
	}
}

func (s *Server) Run() error {
	s.Logger.Info(
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
	s.Logger.Info("Shutting down Rest API server")
	return s.server.Shutdown(ctx)
}

func (s *Server) Name() string {
	return "rest"
}

// Enforces interface definition
var _ app.Component = (*Server)(nil)
