package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/GoodCodingFriends/gpay/adapter"
	"github.com/GoodCodingFriends/gpay/adapter/controller"
	"github.com/GoodCodingFriends/gpay/adapter/repository"
	"github.com/GoodCodingFriends/gpay/adapter/store"
	"github.com/GoodCodingFriends/gpay/config"
)

func main() {
	cfg, err := config.Process()
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "[gpay] ", log.Lshortfile|log.LstdFlags)
	repo, err := repository.NewMySQLRepository(cfg)
	if err != nil {
		panic(err)
	}
	defer repo.Close()

	store, err := store.NewGCSStore(cfg)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	go func() {
		bot, err := controller.NewSlackBot(logger, cfg, repo, store)
		if err != nil {
			panic(err)
		}
		os.Exit(run(bot))
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	<-sigChan
	return
}

func run(listener adapter.Listener) int {
	if err := listener.Listen(); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}
