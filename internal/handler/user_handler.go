package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nazibul7/inmemory-user-api/internal/model"
	"github.com/nazibul7/inmemory-user-api/internal/store"
)

type UserHandler struct {
	store *store.UserStore
}

func NewUserHandler(store *store.UserStore) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

// POST user
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if user.Name == "" || user.Email == "" || user.ID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	if err := h.store.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// GET users
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users := h.store.GetAllUser()
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GET user
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}
	user, err := h.store.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// PUT user/id
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if user.Name == "" || user.Email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	if err := h.store.UpdateUser(id, user); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// DELETE user/id
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
