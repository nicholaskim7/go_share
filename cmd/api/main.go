package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nicholaskim7/go_share/internal/handlers"
	"github.com/nicholaskim7/go_share/internal/middleware"
	"github.com/nicholaskim7/go_share/internal/services"
	"github.com/nicholaskim7/go_share/internal/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// database set up
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Database connection established")

	// dependencies
	postStore := storage.NewPostDBStore(db)
	userStore := storage.NewUserDBStore(db)
	userService := services.NewUserService(userStore)
	// inject Store/service dependency to Handlers
	postHandler := handlers.NewPostHandler(postStore)
	userHandler := handlers.NewUserHandler(userService)

	// public routes (no auth needed)
	http.HandleFunc("GET /posts", postHandler.GetPosts)
	http.HandleFunc("GET /posts/user/{username}", postHandler.GetPostsByUsername)
	http.HandleFunc("GET /posts/tag/{tag}", postHandler.GetPostsByTag)
	// if needed refractor user routing like posts
	http.Handle("/users", userHandler)
	http.HandleFunc("/login", userHandler.SignIn)

	// protected routes wrapped in middleware
	http.HandleFunc("POST /posts", middleware.AuthMiddleware(postHandler.CreatePost))

	addr := ":8080"
	server := &http.Server{
		Addr:         addr,
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Printf("Server starting on http://localhost%s\n", addr)
	log.Fatal(server.ListenAndServe())
}
