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
	if err := godotenv.Load(); err != nil {
		// in Docker, variables come from the environment not a file
		log.Println("No .env file found, relying on system environment variables")
	}

	// database set up with Retry
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Printf("sql.Open failed: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		err = db.Ping()
		if err == nil {
			log.Println("Database connection established")
			break
		}
		log.Printf("Failed to connect to database (attempt %d/10): %v", i+1, err)
		log.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to database after retries:", err)
	}
	defer db.Close()

	// create uploads directory
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", 0755)
	}

	// dependencies
	postStore := storage.NewPostDBStore(db)
	userStore := storage.NewUserDBStore(db)
	userService := services.NewUserService(userStore)
	// inject Store/service dependency to Handlers
	postHandler := handlers.NewPostHandler(postStore)
	userHandler := handlers.NewUserHandler(userService, userStore)

	// public routes (no auth needed)
	http.HandleFunc("GET /posts", postHandler.GetPosts)
	http.HandleFunc("GET /posts/user/{username}", postHandler.GetPostsByUsername)
	http.HandleFunc("GET /posts/tag/{tag}", postHandler.GetPostsByTag)

	http.HandleFunc("POST /users", userHandler.CreateUser)
	http.HandleFunc("GET /users", userHandler.GetUsers)
	http.HandleFunc("POST /login", userHandler.SignIn)
	http.HandleFunc("GET /users/user/{username}", userHandler.GetUserByUsername)

	// protected routes wrapped in middleware
	http.HandleFunc("POST /posts", middleware.AuthMiddleware(postHandler.CreatePost))
	// cannot logout if not logged in
	http.HandleFunc("POST /logout", middleware.AuthMiddleware(userHandler.SignOut))
	http.HandleFunc("DELETE /posts/id/{id}", middleware.AuthMiddleware(postHandler.DeletePostById))

	// serve static files
	fs := http.FileServer(http.Dir("./uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads", fs))

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
