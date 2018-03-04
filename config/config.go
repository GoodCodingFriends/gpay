package config

import (
	"sync"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/kelseyhightower/envconfig"
)

var (
	once sync.Once
	conf Config
)

type Config struct {
	entity *entity.Config
}

var processErr error

func Process() (*Config, error) {
	once.Do(func() {
		processErr = envconfig.Process("", &conf)
	})
	return &conf, processErr
}
