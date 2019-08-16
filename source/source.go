package source

import (
	"context"
	"io"
)

// Source represents image source.
type Source interface {
	// Random returns an image handle randomly choiced.
	Random(ctx context.Context) (io.ReadCloser, error)
}
