package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nazibul7/inmemory-user-api/internal/model"
	"github.com/nazibul7/inmemory-user-api/internal/store"
)

func setup() *UserHandler {
	store := store.NewUserStore()
	handler := NewUserHandler(store)
	return handler
}

// --------- Create----------------------------------------------------------
func TestCreate_Success(t *testing.T) {
	user := model.User{
		ID:    "1",
		Name:  "Nazibul",
		Email: "nazibul@example.com",
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal user %v", err)
	}
	handler := setup()
	// NewBuffer used instead of bytes.Buffer because we now already no more data will come & it is constant
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Create(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", res.StatusCode)
	}

	var created model.User
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if created.ID != user.ID {
		t.Errorf("expected ID %s, got %s", user.ID, created.ID)
	}

}

func TestCreate_InvalidJSON(t *testing.T) {
	handler := setup()

	invalidJSON := []byte(`{"id": "1", "name":`)

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Create(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", res.StatusCode)
	}
}
func TestCreate_MissingFields(t *testing.T) {
	user := model.User{ID: "1"}
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal user %v", err)
	}
	handler := setup()
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Create(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", res.StatusCode)
	}
}

func TestCreate_Duplicate(t *testing.T) {
	user := model.User{
		ID:    "1",
		Name:  "Nazibul",
		Email: "nazibul@example.com",
	}
	handler := setup()
	// already storing a user in store
	handler.store.CreateUser(user)

	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal user %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Create(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusConflict {
		t.Errorf("expected 409, got %d", res.StatusCode)
	}
}

//------------GetAll-----------------------------------------------

func TestGetAll_Empty(t *testing.T) {
	handler := setup()

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	res := w.Result()

	var users []model.User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("expected 0 got %v", len(users))
	}
}

func TestGetAll_ReturnsAll(t *testing.T) {
	handler := setup()
	user1 := model.User{
		ID:    "1",
		Name:  "Nazibul",
		Email: "nazibul@example.com",
	}
	user2 := model.User{
		ID:    "2",
		Name:  "Hossain",
		Email: "hossain@example.com",
	}
	handler.store.CreateUser(user1)
	handler.store.CreateUser(user2)

	req := httptest.NewRequest(http.MethodGet, "/user", nil)

	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	res := w.Result()

	var users []model.User

	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

//--------GetUser--------------------------------------------------

func TestGetUser_Success(t *testing.T) {
	handler := setup()
	user := model.User{
		ID:    "1",
		Name:  "Nazibul",
		Email: "nazibul@example.com",
	}
	handler.store.CreateUser(user)
	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	res := w.Result()

	var userData model.User
	if err := json.NewDecoder(res.Body).Decode(&userData); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if userData.ID != user.ID {
		t.Errorf("expected ID %s, got %s", user.ID, userData.ID)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	handler := setup()

	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected not found 404, got %d", res.StatusCode)
	}
}

func TestGetUser_MissingID(t *testing.T) {
	handler := setup()

	req := httptest.NewRequest(http.MethodGet, "/user/", nil)
	req.SetPathValue("id", "")

	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected bad request 400, got %d", res.StatusCode)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	handler := setup()
	user := model.User{
		ID:    "1",
		Name:  "Nazibul",
		Email: "nazibul@example.com",
	}

	handler.store.CreateUser(user)

	updated := model.User{ID: "1", Name: "Updated Name", Email: "updated@example.com"}

	body, err := json.Marshal(updated)
	if err != nil {
		t.Fatalf("failed to marshal user %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/user/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", res.StatusCode)
	}
	var resUser model.User
	if err := json.NewDecoder(res.Body).Decode(&resUser); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if resUser.Name != updated.Name {
		t.Errorf("expected name %s, got %s", updated.Name, resUser.Name)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	handler := setup()
	user := model.User{
		ID:    "1",
		Name:  "Nazibul",
		Email: "nazibul@example.com",
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal user %v", err)
	}
	req := httptest.NewRequest(http.MethodPut, "/user/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.UpdateUser(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", res.StatusCode)
	}
}

func TestUpdateUser_MissingFields(t *testing.T) {
	handler := setup()
	user := model.User{
		ID:    "1",
		Name:  "Nazibul",
		Email: "nazibul@example.com",
	}
	handler.store.CreateUser(user)

	updated := model.User{ID: "1"} // missing Name & Email
	body, err := json.Marshal(updated)
	if err != nil {
		t.Fatalf("failed to marshal user %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/user/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.UpdateUser(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	handler := setup()
	user := model.User{
		ID:    "1",
		Name:  "Nazibul",
		Email: "nazibul@example.com",
	}
	handler.store.CreateUser(user)

	req := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.DeleteUser(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", res.StatusCode)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	handler := setup()
	req := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()

	handler.DeleteUser(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", res.StatusCode)
	}
}
