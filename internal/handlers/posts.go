package handlers

// implementing handlers that will use db functions
import (
	"encoding/json"
	"net/http"

	"github.com/nicholaskim7/go_share/internal/models"
	"github.com/nicholaskim7/go_share/internal/storage"
)

type PostHandler struct {
	store storage.PostStore
}

// ensure that every new post handler has a store specifically Poststore
func NewPostHandler(store storage.PostStore) *PostHandler {
	return &PostHandler{store: store}
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getPosts(w, r)
	case http.MethodPost:
		h.createPost(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *PostHandler) getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	posts, err := h.store.GetAll(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch posts", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *PostHandler) createPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var newPost models.Post
	// decode request body into new post
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	// minimal validation post must have at least title and body
	if newPost.Title == "" || newPost.Body == "" {
		http.Error(w, "title or body is required", http.StatusBadRequest)
		return
	}
	// call db method create to insert new post
	created, err := h.store.Create(r.Context(), newPost)
	if err != nil {
		http.Error(w, "failed to create post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}
