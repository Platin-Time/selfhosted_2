package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"ci-demo/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /stats", handlers.Stats)
	mux.HandleFunc("POST /calc", handlers.Calc)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}
