package handlers

// implementing handlers that will use db functions
import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/nicholaskim7/go_share/internal/middleware"
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
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "posts not found", http.StatusNotFound) // 404
			return
		}
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
	// retrieve userID from context set by middleware
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "internal auth error", http.StatusInternalServerError)
		return
	}
	// must use multipart/form-data to send files along with text
	err := r.ParseMultipartForm(10 << 20) //limit file size to 10MB
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}
	// extract the normal post data
	var newPost models.Post
	newPost.Title = r.FormValue("title")
	newPost.Body = r.FormValue("body")
	if values, ok := r.MultipartForm.Value["tags"]; ok {
		newPost.Tags = values
	}
	// force author id to be logged in user
	newPost.UserID = userID
	// minimal validation post must have at least title and body
	if newPost.Title == "" || newPost.Body == "" {
		http.Error(w, "title or body is required", http.StatusBadRequest)
		return
	}
	// handle file uploads
	formFiles := r.MultipartForm.File["files"]
	for _, fileHeader := range formFiles {
		// open the file stream
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "error reading file", http.StatusInternalServerError)
			return
		}
		// generate unique name (timestamp + file name)
		uniqueFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
		// create the file on the servers disk
		// in the future replace with s3 storage
		dst, err := os.Create("./uploads/" + uniqueFileName)
		if err != nil {
			file.Close()
			http.Error(w, "server storage error", http.StatusInternalServerError)
			return
		}
		// copy content to the destination file
		_, err = io.Copy(dst, file)
		//close both files
		dst.Close()
		file.Close()
		if err != nil {
			http.Error(w, "error saving file", http.StatusInternalServerError)
			return
		}
		// append the stored filename to the struct
		newPost.Files = append(newPost.Files, uniqueFileName)
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
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "posts from username not found", http.StatusNotFound) // 404
			return
		}
		http.Error(w, "failed to fetch posts by username", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetPostsByTag(w http.ResponseWriter, r *http.Request) {
	// get tag from url path
	// maintain restfullness of get request
	tag := r.PathValue("tag")
	if tag == "" {
		http.Error(w, "no tag provided", http.StatusBadRequest)
		return
	}
	posts, err := h.store.GetByTag(r.Context(), tag)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "posts by tag not found", http.StatusNotFound) // 404
			return
		}
		http.Error(w, "failed to fetch posts by tag", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
