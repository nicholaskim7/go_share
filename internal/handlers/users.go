package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nicholaskim7/go_share/internal/models"
	"github.com/nicholaskim7/go_share/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
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
	// call db method create to insert new user
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
