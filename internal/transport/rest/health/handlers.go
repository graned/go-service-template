package health

import (
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Get(r *http.Request) (HealthResponse, error) {
	return HealthResponse{
		Status: "Ok",
	}, nil
}
