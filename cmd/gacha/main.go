package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoodCodingFriends/gpay/slack"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/slack", slack.Router)

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("server listen in %s", addr)
	if err := http.ListenAndServe(addr, r); err != http.ErrServerClosed {
		log.Fatalf("failed to launch the server: %s", err)
	}
}
