package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
)

var (
	once sync.Once
	conf Config
)

type Config struct {
	BalanceLowerLimit int64
}

var processErr error

func Process() (*Config, error) {
	once.Do(func() {
		processErr = envconfig.Process("", &conf)
	})
	return &conf, processErr
}
