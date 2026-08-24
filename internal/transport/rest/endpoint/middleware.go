package endpoint

import (
	"net/http"
)

type Middleware func(Handler) Handler

func Chain(
	handler Handler,
	middlewares ...Middleware,
) Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

// This is a sample on how the auth middleware could look like
func RequireAuth(
	next Handler,
) Handler {
	return func(r *http.Request) (Result, error) {
		token := r.Header.Get("Authorization")

		if token == "" {
			return Result{}, ErrUnauthorized
		}

		return next(r)
	}
}
