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

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" || user.ID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// WithCancel — not needed here
	// r.Context() already has cancellation (client disconnect / shutdown)
	// Only use WithTimeout/Cancel when need custom control
	ctx := r.Context()

	if err := h.store.CreateUser(ctx, user); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// WHY Encode here:
	// - Simple response
	// - Small payload
	// - No need for strict control
	// - If encoding fails, very rare and acceptable here
	//
	// Encode writes directly to ResponseWriter:
	// → automatically sends status (201 already set)
	// → streams response

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		// If this fails, headers already sent → can't change status
		// acceptable tradeoff for simple APIs
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// GET users
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.GetAllUser(r.Context())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// WHY Encode here:
	// - Read operation
	// - No status change needed (default 200)
	// - Simpler and cleaner
	// - Good for most REST APIs

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// GET user
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	user, err := h.store.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// WHY Encode:
	// - Small single object
	// - No strict need to pre-buffer
	// - Simpler

	w.Header().Set("Content-Type", "application/json")

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

	if err := h.store.UpdateUser(r.Context(), id, user); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Here we demonstrate Marshal usage (STRICT pattern)

	// WHY Marshal here:
	// - We want to ensure encoding succeeds BEFORE sending response
	// - Avoid partial/broken response
	// - Better for production-critical endpoints

	data, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// DELETE user/id
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// No body needed → just status
	// No Encode / Marshal required

	w.WriteHeader(http.StatusNoContent)
}