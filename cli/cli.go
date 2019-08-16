package cli

import (
	"context"
	"io"
	"os"
	"sync"
	"time"

	"github.com/GoodCodingFriends/gpay/source/gcs"
	"github.com/GoodCodingFriends/gpay/usecase"
	"github.com/morikuni/failure"
)

var (
	InvalidUsageCode      = failure.StringCode("invalid usage")
	UnknownSubCommandCode = failure.StringCode("unknown command")
)

type CLI struct {
	Reader            io.Reader
	Writer, ErrWriter io.Writer

	initOnce sync.Once
}

func (c *CLI) Run(args []string) error {
	if len(args) == 0 {
		return failure.New(InvalidUsageCode)
	}

	switch args[0] {
	case "lgtm":
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		var err error
		c.initOnce.Do(func() {
			c.initCommon()
			if len(args) == 1 {
				return
			}
			src, berr := gcs.New(ctx, args[1:])
			if berr != nil {
				err = berr
				return
			}
			usecase.Inject(
				usecase.InteractorParams{
					Source: src,
				},
			)
		})
		if err != nil {
			return failure.Wrap(err)
		}
		err = usecase.OverlayLGTM(ctx, c.Writer)
		if err != nil {
			return failure.Wrap(err)
		}
	default:
		return failure.New(UnknownSubCommandCode, failure.Context{"name": args[0]})
	}
	return nil
}

func (c *CLI) initCommon() {
	if c.Reader == nil {
		c.Reader = os.Stdin
	}
	if c.Writer == nil {
		c.Writer = os.Stdout
	}
	if c.ErrWriter == nil {
		c.ErrWriter = os.Stderr
	}
}
