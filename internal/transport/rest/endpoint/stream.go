package endpoint

import (
	"context"
	"net/http"
)

type Stream func(
	ctx context.Context,
	w http.ResponseWriter,
) error

type StreamResult struct {
	Status      int
	ContentType string
	Stream      Stream
}
