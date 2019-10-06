package cli

import (
	"context"
	"io"
	"os"
	"sync"

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
	ctx := context.Background()
	var err error
	c.initOnce.Do(func() {
		c.initCommon()

		var lgtm bool
		if len(args) == 0 {
			args = []string{"eupho", "hamachi", "hanairo", "saekano", "kaguya"}
		} else if args[0] == "lgtm" {
			args = args[1:]
			lgtm = true
		}
		src, berr := gcs.New(ctx, args, lgtm)
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
	err = usecase.FetchImage(ctx, c.Writer)
	if err != nil {
		return failure.Wrap(err)
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
