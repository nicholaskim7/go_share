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

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.store.GetAll(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
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

func (h *PostHandler) GetPostsByUsername(w http.ResponseWriter, r *http.Request) {
	// get username from url path
	// maintain restfullness of get request
	userName := r.PathValue("username")
	if userName == "" {
		http.Error(w, "no username provided", http.StatusBadRequest)
		return
	}
	posts, err := h.store.GetByUsername(r.Context(), userName)
	if err != nil {
		http.Error(w, "failed to fetch posts by username", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
