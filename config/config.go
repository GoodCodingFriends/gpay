package config

import (
	"sync"

	"github.com/k0kubun/pp"
	"github.com/kelseyhightower/envconfig"
)

var (
	once sync.Once
	cfg  Config
)

type Config struct {
	Meta       *Meta
	Controller *Controller
	Entity     *Entity
}

type Meta struct {
	Debug bool `default:"true"`
}

type Entity struct {
	BalanceLowerLimit int64 `default:"-5000"`
}

type Controller struct {
	Slack *Slack
}

type Slack struct {
	APIToken string
	BotName  string `default:"gpay"`
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
