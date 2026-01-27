package main

import (
	"log"
	"net/http"

	"github.com/nicholaskim7/go_share/internal/handlers"
	"github.com/nicholaskim7/go_share/internal/storage"
)

func main() {
	postStore := storage.NewPostStore()
	userStore := storage.NewUserStore()
	// inject Store dependency to Handlers
	postHandler := handlers.NewPostHandler(postStore)
	userHandler := handlers.NewUserHandler(userStore)
	http.Handle("/posts", postHandler)
	http.Handle("/users", userHandler)

	addr := ":8080"

	log.Printf("Server starting on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
