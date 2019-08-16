package slack

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/morikuni/failure"
	nslack "github.com/nlopes/slack"
)

type apiClient struct {
	c     *nslack.Client
	token string
}

func newAPIClient() *apiClient {
	token := os.Getenv("SLACK_ACCESS_TOKEN")
	if token == "" {
		panic("$SLACK_ACCESS_TOKEN is misisng")
	}
	return &apiClient{
		c:     nslack.New(token),
		token: token,
	}
}

func (c *apiClient) UploadFile(ctx context.Context, r io.Reader, channel string) error {
	name := fmt.Sprintf("%d.png", time.Now().UnixNano())
	f, err := c.c.UploadFileContext(ctx, nslack.FileUploadParameters{
		Title:    name,
		Filename: name,
		Filetype: "png",
		Reader:   r,
		Channels: []string{channel},
	})
	if err != nil {
		return failure.Wrap(err)
	}
	log.Println(f)
	return nil
}
