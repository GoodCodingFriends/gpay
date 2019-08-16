package usecase

import (
	"context"
	"io"

	"github.com/morikuni/failure"
)

func FetchLGTM(ctx context.Context, w io.Writer) error {
	return defaultInteractor.FetchLGTM(ctx, w)
}
func (i *interactor) FetchLGTM(ctx context.Context, w io.Writer) error {
	r, err := i.Source.Random(ctx)
	if err != nil {
		return failure.Wrap(err)
	}
	defer r.Close()
	if _, err := io.Copy(w, r); err != nil {
		return failure.Wrap(err)
	}
	return nil
}
