package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GoodCodingFriends/gpay/adapter"
	"github.com/GoodCodingFriends/gpay/adapter/controller"
	"github.com/GoodCodingFriends/gpay/config"
)

func main() {
	cfg, err := config.Process()
	if err != nil {
		panic(err)
	}
	logger := log.New(os.Stdout, "[gpay] ", log.Lshortfile|log.LstdFlags)

	os.Exit(run(controller.NewSlackBot(logger, cfg)))
}

func run(listener adapter.Listener) int {
	if err := listener.Listen(); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}
