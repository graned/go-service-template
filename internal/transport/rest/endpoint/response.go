// Typed of standard API response
package endpoint

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Status int            `json:"status"`
	Data   any            `json:"data"`
	Error  *ErrorResponse `json:"error"`
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	response Envelope,
) error {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(response)
}

func writeSuccess(
	w http.ResponseWriter,
	result Result,
) error {
	return writeJSON(
		w,
		result.Status,
		Envelope{
			Status: result.Status,
			Data:   result.Data,
			Error:  nil,
		},
	)
}

func WriteError(
	w http.ResponseWriter,
	err error,
) error {
	status, apiError := mapError(err)

	return writeJSON(
		w,
		status,
		Envelope{
			Status: status,
			Data:   nil,
			Error:  &apiError,
		},
	)
}
