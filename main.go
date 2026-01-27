package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Post struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Tags        []string  `json:"tags"`
	Files       []string  `json:"files"`
	DateCreated time.Time `json:"date_created"`
}

var posts = []Post{
	{ID: 1, UserID: 1, Title: "python script", Body: "code contains script", Tags: []string{"programming", "python", "script"}, Files: []string{"script.py"}, DateCreated: time.Now().UTC()},
	{ID: 2, UserID: 1, Title: "go script", Body: "code contains calculator app", Tags: []string{"coding", "python", "go"}, Files: []string{"main.go"}, DateCreated: time.Now().UTC()},
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPosts(w, r)
	case http.MethodPost:
		createPost(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var newPost Post
	var postsMu sync.Mutex
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	newPost.DateCreated = time.Now().UTC()
	postsMu.Lock()
	posts = append(posts, newPost)
	postsMu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newPost); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	addr := ":8080"
	http.HandleFunc("/posts", postHandler)
	log.Printf("Server starting on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
