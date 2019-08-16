package usecase

import (
	"context"
	"io"

	"github.com/ktr0731/lgtm"
	"github.com/morikuni/failure"
)

func OverlayLGTM(ctx context.Context, w io.Writer) error {
	return defaultInteractor.OverlayLGTM(ctx, w)
}
func (i *interactor) OverlayLGTM(ctx context.Context, w io.Writer) error {
	r, err := i.Source.Random(ctx)
	if err != nil {
		return failure.Wrap(err)
	}
	defer r.Close()
	if err := lgtm.New(r, w); err != nil {
		return failure.Wrap(err)
	}
	return nil
}
