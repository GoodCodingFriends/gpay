package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoodCodingFriends/gpay/slack"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/slack", slack.Handler())
	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if err := http.ListenAndServe(addr, nil); err != http.ErrServerClosed {
		log.Fatalf("failed to launch the server: %s", err)
	}
}
