package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/nazibul7/inmemory-user-api/internal/middleware"
	"github.com/nazibul7/inmemory-user-api/internal/model"
)

type UserStorer interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, id string) (model.User, error)
	GetAllUser(ctx context.Context) ([]model.User, error)
	UpdateUser(ctx context.Context, id string, user model.User) error
	DeleteUser(ctx context.Context, id string) error
}

type UserHandler struct {
	store UserStorer
}

func NewUserHandler(store UserStorer) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func getRequestID(r *http.Request) string {
	if id, ok := r.Context().Value(middleware.RequestIDKey).(string); ok {
		return id
	}
	return "no-request-id"
}

// POST user
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	reqID := getRequestID(r)
	log.Printf("requestID=%s Create called", reqID)

	var user model.User

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("requestID=%s Create invalid JSON: %v", reqID, err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" || user.ID == "" {
		log.Printf("requestID=%s Create missing required fields", reqID)
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// WithCancel — not needed here
	// r.Context() already has cancellation (client disconnect / shutdown)
	// Only use WithTimeout/Cancel when need custom control
	ctx := r.Context()

	if err := h.store.CreateUser(ctx, user); err != nil {
		log.Printf("requestID=%s conflict creating user id=%s: %v", reqID, user.ID, err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	log.Printf("requestID=%s user created id=%s", reqID, user.ID)

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
	reqID := getRequestID(r)
	log.Printf("requestID=%s GetAll called", reqID)

	users, err := h.store.GetAllUser(r.Context())
	if err != nil {
		log.Printf("requestID=%s GetAll error fetching users: %v", reqID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("requestID=%s GetAll returning %d users", reqID, len(users))

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
	reqID := getRequestID(r)
	log.Printf("requestID=%s GetUser called", reqID)

	id := r.PathValue("id")
	if id == "" {
		log.Printf("requestID=%s missing id in GetUser", reqID)
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	user, err := h.store.GetUser(r.Context(), id)
	if err != nil {
		log.Printf("requestID=%s user not found id=%s in GetUser: %v", reqID, id, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("requestID=%s user found id=%s", reqID, id)

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
	reqID := getRequestID(r)
	log.Printf("requestID=%s UpdateUser called", reqID)

	id := r.PathValue("id")
	if id == "" {
		log.Printf("requestID=%s missing id in UpdateUser", reqID)
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("requestID=%s UpdateUser invalid JSON: %v", reqID, err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" {
		log.Printf("requestID=%s UpdateUser missing required fields", reqID)
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if err := h.store.UpdateUser(r.Context(), id, user); err != nil {
		log.Printf("requestID=%s user not found id=%s in UpdateUser: %v", reqID, id, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("requestID=%s user updated id=%s", reqID, id)

	// Here we demonstrate Marshal usage (STRICT pattern)

	// WHY Marshal here:
	// - We want to ensure encoding succeeds BEFORE sending response
	// - Avoid partial/broken response
	// - Better for production-critical endpoints

	data, err := json.Marshal(user)
	if err != nil {
		log.Printf("requestID=%s UpdateUser marshal error: %v", reqID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// DELETE user/id
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	reqID := getRequestID(r)
	log.Printf("requestID=%s DeleteUser called", reqID)

	id := r.PathValue("id")
	if id == "" {
		log.Printf("requestID=%s missing id in DeleteUser", reqID)
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteUser(r.Context(), id); err != nil {
		log.Printf("requestID=%s user not found id=%s in DeleteUser: %v", reqID, id, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("requestID=%s user deleted id=%s", reqID, id)

	// No body needed → just status
	// No Encode / Marshal required

	w.WriteHeader(http.StatusNoContent)
}