package main

import (
	"fmt"
	"github.com/hayashiki/tarsier-integration/handler"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	h := handler.NewHandler()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/slack/invoke", func(w http.ResponseWriter, r *http.Request) {
		h.InvokeSlackAuth(w, r)
	})
	r.Get("/slack/callback", func(w http.ResponseWriter, r *http.Request) {
		h.HandleSlackAuth(w, r)
		//handler.Wrap(h.HandleSlackAuth).ServeHTTP(w, r)
	})
	r.Post("/slack/interactive", func(w http.ResponseWriter, r *http.Request) {
		handler.Wrap(h.HandleSlackInteractive).ServeHTTP(w, r)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
