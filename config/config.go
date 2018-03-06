package config

import (
	"sync"

	"github.com/GoodCodingFriends/gpay/adapter/controller"
	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/k0kubun/pp"
	"github.com/kelseyhightower/envconfig"
)

var (
	once sync.Once
	cfg  Config
)

type Config struct {
	Meta       *Meta
	Controller *controller.Config
	Entity     *entity.Config
}

type Meta struct {
	Debug bool `default:"true"`
}

var processErr error

func Process() (*Config, error) {
	once.Do(func() {
		processErr = envconfig.Process("", &cfg)
	})
	if cfg.Meta.Debug {
		pp.Println(cfg)
	}
	return &cfg, processErr
}
