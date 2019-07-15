package slack

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var signingSecret []byte

// TODO: use config.
func init() {
	s := os.Getenv("SLACK_SIGNING_SECRET")
	if s == "" {
		panic("$SLACK_SIGNING_SECRET is required")
	}
	signingSecret = []byte(s)
}

func authenticationMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts := r.Header.Get("X-Slack-Request-Timestamp")
		if ts == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		its, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			// TODO: logging.
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		t := time.Unix(its, 0)
		if t.Before(time.Now().Add(-5 * time.Minute)) {
			// TODO: logging.
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var buf bytes.Buffer
		io.Copy(&buf, r.Body)
		r.Body.Close()
		r.Body = ioutil.NopCloser(&buf)

		hash := hmac.New(sha256.New, signingSecret)
		fmt.Fprintf(hash, "v0:%s:%s", ts, buf.String())
		sig := fmt.Sprintf("v0=%x", hash.Sum(nil))
		if sig != r.Header.Get("X-Slack-Signature") {
			log.Printf("expected signature: %s, but got %s", sig, r.Header.Get("X-Slack-Signature"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
