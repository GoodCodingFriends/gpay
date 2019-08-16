package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/morikuni/failure"
)

type apiClient struct {
	c     *http.Client
	token string
}

func newAPIClient() *apiClient {
	token := os.Getenv("SLACK_ACCESS_TOKEN")
	if token == "" {
		panic("$SLACK_ACCESS_TOKEN is misisng")
	}
	return &apiClient{
		c:     http.DefaultClient,
		token: token,
	}
}

func (c *apiClient) UploadFile(ctx context.Context, r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return failure.Wrap(err)
	}
	var in bytes.Buffer
	enc := json.NewEncoder(&in)
	err = enc.Encode(struct {
		Token   string `json:"token"`
		Content []byte `json:"content"`
	}{
		Token:   c.token,
		Content: b,
	})
	if err != nil {
		return failure.Wrap(err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://slack.com/api/files.upload", &in)
	if err != nil {
		return failure.Wrap(err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := c.c.Do(req)
	if err != nil {
		return failure.Wrap(err)
	}
	defer res.Body.Close()

	return nil
}
