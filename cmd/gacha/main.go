package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoodCodingFriends/gpay/slack"
	"github.com/go-chi/chi"
)

func main() {
	mux := chi.NewRouter()
	mux.Route("/slack", slack.Router)
	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("server listen in %s", addr)
	if err := http.ListenAndServe(addr, mux); err != http.ErrServerClosed {
		log.Fatalf("failed to launch the server: %s", err)
	}
}
