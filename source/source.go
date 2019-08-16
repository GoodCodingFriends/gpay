package source

import (
	"context"
	"io"

	"github.com/morikuni/failure"
)

var InvalidParameterCode = failure.StringCode("invalid param")

// Source represents image source.
type Source interface {
	// Random returns an image handle randomly choiced.
	Random(ctx context.Context) (io.ReadCloser, error)
}
