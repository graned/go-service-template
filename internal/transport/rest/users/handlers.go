package users

import (
	applogger "github.com/graned/go-service-template/internal/logger"
	domainuser "github.com/graned/go-service-template/internal/user"
	"net/http"
)

type Handler struct {
	service *domainuser.Service
}

func NewHandler(
	service *domainuser.Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Get(r *http.Request) (UserResponse, error) {
	id := r.PathValue("id")
	return UserResponse{
		ID:        id,
		FirstName: "Eduardo",
		LastName:  "Anaya",
		Email:     "eduardo@example.com",
	}, nil
}

func (h *Handler) List(r *http.Request) ([]UserResponse, error) {
	var res []UserResponse

	res = append(res, UserResponse{
		ID:        "123",
		FirstName: "Eduardo",
		LastName:  "Anaya",
		Email:     "eduardo@example.com",
	})
	return res, nil
}

func (h *Handler) Create(r *http.Request, body CreateUserRequest) (UserResponse, error) {
	logger := applogger.FromContext(r.Context())

	logger.InfoContext(
		r.Context(),
		"creating user",
	)

	return UserResponse{
		ID:        "1234",
		FirstName: "Eduardo",
		LastName:  "Anaya",
		Email:     "eduardo@example.com",
	}, nil
}
