package controller

type Config struct {
	Slack *SlackConfig
}

type SlackConfig struct {
	APIToken string
	BotName  string `default:"gpay"`
}
