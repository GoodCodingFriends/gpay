package main

import (
	"fmt"
	"os"

	"github.com/GoodCodingFriends/gpay/adapter"
	"github.com/GoodCodingFriends/gpay/adapter/controller"
	"github.com/GoodCodingFriends/gpay/config"
	"github.com/k0kubun/pp"
)

func main() {
	cfg, err := config.Process()
	if err != nil {
		panic(err)
	}
	pp.Println(cfg)
	os.Exit(run(controller.NewSlackBot(cfg.Controller.Slack)))
}

func run(listener adapter.Listener) int {
	if err := listener.Listen(); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}
