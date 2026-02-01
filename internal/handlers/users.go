package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nicholaskim7/go_share/internal/auth"
	"github.com/nicholaskim7/go_share/internal/models"
	"github.com/nicholaskim7/go_share/internal/services"
	"github.com/nicholaskim7/go_share/internal/storage"
)

type UserHandler struct {
	service *services.UserService
	store   *storage.UserDBStore
}

func NewUserHandler(service *services.UserService, store *storage.UserDBStore) *UserHandler {
	return &UserHandler{service: service, store: store}
}

// func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		h.getUsers(w, r)
// 	case http.MethodPost:
// 		h.createUser(w, r)
// 	default:
// 		w.Header().Set("Allow", "GET, POST")
// 		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch users", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var newUser models.User
	// decode request body into new user
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	// minimal validation, user must have required fields
	if newUser.FirstName == "" || newUser.LastName == "" || newUser.UserName == "" || newUser.Email == "" || newUser.Password == "" {
		http.Error(w, "fill in required fields", http.StatusBadRequest)
		return
	}
	// call user service which will call db method create to insert new user after password hash
	created, err := h.service.CreateUser(r.Context(), newUser)
	if err != nil {
		http.Error(w, "failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// clear password hash so it doesnt get sent back to client
	created.Password = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var userLoginPayload models.UserLoginPayload
	if err := json.NewDecoder(r.Body).Decode(&userLoginPayload); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	if userLoginPayload.UserName == "" || userLoginPayload.Password == "" {
		http.Error(w, "fill in required credentials", http.StatusBadRequest)
		return
	}
	// call login service
	user, err := h.service.LoginUser(r.Context(), userLoginPayload)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	// generate JWT token
	token, err := auth.CreateToken(user.ID)
	if err != nil {
		http.Error(w, "error generating session", http.StatusInternalServerError)
		return
	}
	// set http-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,  // javascript cannot read this (No XSS)
		Secure:   false, // set to true in production (requires https)
		Path:     "/",
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		http.Error(w, "no username provided", http.StatusBadRequest)
		return
	}
	user, err := h.store.GetByUsername(r.Context(), username)
	if err != nil {
		http.Error(w, "failed to fetch user by username", http.StatusInternalServerError)
		return
	}
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
