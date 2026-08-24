package endpoint

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func JSONBody[Request any, Response any](
	validate *validator.Validate,
	status int,
	handler BodyHandler[Request, Response],
) Handler {
	return func(r *http.Request) (Result, error) {
		var body Request

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&body); err != nil {
			return Result{}, ApiBadRequest(
				"invalid_json",
				"invalid request body",
				err,
			)
		}

		if err := validate.Struct(body); err != nil {
			return Result{}, ApiValidationError(err)
		}

		data, err := handler(
			r,
			body,
		)

		if err != nil {
			return Result{}, err
		}

		return Result{
			Status: status,
			Data:   data,
		}, nil
	}
}
