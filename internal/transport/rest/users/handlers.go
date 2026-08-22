package users

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	applogger "github.com/graned/go-service-template/internal/logger"
	domainuser "github.com/graned/go-service-template/internal/user"
	"net/http"
)

type Handler struct {
	service   *domainuser.Service
	validator *validator.Validate
}

func NewHandler(
	service *domainuser.Service,
	validator *validator.Validate,
) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
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
	var request CreateUserRequest

	logger := applogger.FromContext(r.Context())

	logger.InfoContext(
		r.Context(),
		"creating user",
	)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
