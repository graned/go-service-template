package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/graned/go-service-template/internal/transport"
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
	r.Use(middleware.Recoverer)

	validate := validator.New(
		validator.WithRequiredStructEnabled(),
	)
	usersHandlers := users.NewHandler(userService, validate)

	r.Get("/health", health)
	r.Mount("/users", users.Routes(usersHandlers))

	return &Server{
		Runtime: runtime,

		server: &http.Server{
			Addr:    runtime.Address,
			Handler: r,
		},
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Ok"))
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

// Enforces interface definition
var _ transport.Server = (*Server)(nil)
