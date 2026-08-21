package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/graned/go-service-template/internal/transport/rest/users"
	domainuser "github.com/graned/go-service-template/internal/user"
)

type Server struct {
	server *http.Server
}

func New(
	address string,
	userService *domainuser.Service,
) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	usersHandlers := users.NewHandler(userService)

	r.Get("/health", health)
	r.Mount("/users", users.Routes(usersHandlers))

	return &Server{
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
	fmt.Printf("Rest API Server listening on %s\n", s.server.Addr)

	err := s.server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
