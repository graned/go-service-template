package users

import (
	"fmt"
	domainuser "github.com/graned/go-service-template/internal/user"
	"net/http"
)

type Handler struct {
	service *domainuser.Service
}

func NewHandler(service *domainuser.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprint(w, "user: %s", id)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "users")
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}
