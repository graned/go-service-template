// Centrilized HTTP error representation
package endpoint

import (
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type APIError struct {
	Status  int
	Code    string
	Message string
	Cause   error
}

func (e *APIError) Error() string {
	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Cause
}

func ApiBadRequest(
	code string,
	message string,
	cause error,
) error {
	return &APIError{
		Status:  http.StatusBadRequest,
		Code:    "bad_request",
		Message: message,
		Cause:   cause,
	}
}

func ApiValidationError(
	cause error,
) error {
	return &APIError{
		Status:  http.StatusBadRequest,
		Code:    "validation_failed",
		Message: "Request validation failed",
		Cause:   cause,
	}
}

func ApiMethodNotAllowedError(
	cause error,
) error {
	return &APIError{
		Status:  http.StatusMethodNotAllowed,
		Code:    "method_not_allowed",
		Message: "Method not allowed",
		Cause:   cause,
	}
}

func ApiRouteNotFound(
	cause error,
) error {
	return &APIError{
		Status:  http.StatusNotFound,
		Code:    "route_not_found",
		Message: "Route not found",
		Cause:   cause,
	}
}

func ApiInternalServerError(
	cause error,
) error {
	return &APIError{
		Status:  http.StatusInternalServerError,
		Code:    "internal_server_error",
		Message: "An unexpected error ocurred",
		Cause:   cause,
	}
}

var (
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrInternal     = errors.New("internal server error")
)

func mapError(err error) (int, ErrorResponse) {
	var apiError *APIError

	if errors.As(err, &apiError) {
		return apiError.Status, ErrorResponse{
			Code:    apiError.Code,
			Message: apiError.Message,
		}
	}
	switch {
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest, ErrorResponse{
			Code:    "bad_request",
			Message: "Invalid request",
		}

	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, ErrorResponse{
			Code:    "unauthorized",
			Message: "Unauthorized",
		}

	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden, ErrorResponse{
			Code:    "forbidden",
			Message: "Forbidden",
		}

	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, ErrorResponse{
			Code:    "not_found",
			Message: "Resource not found",
		}

	default:
		return http.StatusInternalServerError, ErrorResponse{
			Code:    "internal_error",
			Message: "An unexpected error occurred",
		}
	}
}
