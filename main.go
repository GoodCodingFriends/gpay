package main

import (
	"os"

	"github.com/GoodCodingFriends/gpay/adapter"
	"github.com/GoodCodingFriends/gpay/adapter/controller"
)

func main() {
	os.Exit(run(controller.NewSlackBot()))
}

func run(listener adapter.Listener) int {
	listener.Listen()
	return 0
}
