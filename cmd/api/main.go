package main

import (
	"log"
	"net/http"

	"github.com/nicholaskim7/go_share/internal/handlers"
	"github.com/nicholaskim7/go_share/internal/storage"
)

func main() {
	store := storage.NewPostStore()
	// inject store dependency to postHandler
	postHandler := handlers.NewPostHandler(store)
	http.Handle("/posts", postHandler)

	addr := ":8080"

	log.Printf("Server starting on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
