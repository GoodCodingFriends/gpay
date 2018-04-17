package store

import (
	"testing"
	"time"

	"github.com/GoodCodingFriends/gpay/config"
	"github.com/stretchr/testify/assert"
)

var cfg *config.Config

func init() {
	var err error
	cfg, err = config.Process()
	if err != nil {
		panic(err)
	}
}

func Test_expired(t *testing.T) {
	d := cfg.Store.GCS.CacheDuration

	cases := []struct {
		d       time.Duration
		expired bool
	}{
		{-d, true},
		{-d + 1*time.Second, false},
	}

	for _, c := range cases {
		cachedAt := time.Now().Add(c.d)
		res := expired(cachedAt, d)
		assert.Equal(t, c.expired, res)
	}
}
