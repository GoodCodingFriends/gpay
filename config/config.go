package config

import (
	"sync"
	"time"

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
	Repository *Repository
	Entity     *Entity
	Store      *Store
}

type Meta struct {
	Debug bool `default:"true"`
}

type Entity struct {
	BalanceLowerLimit int64 `default:"-5000"`
}

type Repository struct {
	MySQL *MySQL
}

type MySQL struct {
	UserName string `default:"root"`
	Password string `default:""`
	Net      string `default:"tcp"`
	Address  string `default:"127.0.0.1:3306"`
	DBName   string `default:"gpay"`
}

type Controller struct {
	Slack *Slack
}

type Slack struct {
	APIToken          string
	DisplayName       string `default:"gPAY"`
	BotName           string `default:"gpay"`
	DoneEmoji         string `default:"sushi"`
	VerificationToken string `default:"sound!euphonium"`

	MaxListTransactionNum int `default:"5"`

	Port string `default:"8080"`
}

type Store struct {
	GCS *GCS
}

type GCS struct {
	CacheDuration time.Duration `default:"12h"`
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
