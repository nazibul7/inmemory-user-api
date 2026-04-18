package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nazibul7/inmemory-user-api/internal/model"
)

type mockUserStore struct {
	CreateFn func(ctx context.Context, user model.User) error
	GetFn    func(ctx context.Context, id string) (model.User, error)
	GetAllFn func(ctx context.Context) ([]model.User, error)
	UpdateFn func(ctx context.Context, id string, user model.User) error
	DeleteFn func(ctx context.Context, id string) error
}

func (m *mockUserStore) CreateUser(ctx context.Context, user model.User) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, user)
	}
	return nil
}
func (m *mockUserStore) GetUser(ctx context.Context, id string) (model.User, error) {
	if m.GetFn != nil {
		return m.GetFn(ctx, id)
	}
	return model.User{}, nil
}
func (m *mockUserStore) GetAllUser(ctx context.Context) ([]model.User, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn(ctx)
	}
	return []model.User{}, nil
}
func (m *mockUserStore) UpdateUser(ctx context.Context, id string, user model.User) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, user)
	}
	return nil
}
func (m *mockUserStore) DeleteUser(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func setup(store *mockUserStore) *UserHandler {
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
	handler := setup(&mockUserStore{})
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
	handler := setup(&mockUserStore{})

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
	handler := setup(&mockUserStore{})
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
	mock := &mockUserStore{CreateFn: func(ctx context.Context, user model.User) error {
		return errors.New("user already exist")
	}}
	handler := setup(mock)

	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal user %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Should use when real db is there maybe in integration test
	// ctx := req.Context()
	// // already storing a user in store
	// handler.store.CreateUser(ctx, user)
	w := httptest.NewRecorder()

	handler.Create(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusConflict {
		t.Errorf("expected 409, got %d", res.StatusCode)
	}
}

//------------GetAll-----------------------------------------------

func TestGetAll_Empty(t *testing.T) {
	handler := setup(&mockUserStore{})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
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
	handler := setup(&mockUserStore{
		GetAllFn: func(ctx context.Context) ([]model.User, error) {
			return []model.User{
				{ID: "1", Name: "Nazibul", Email: "nazibul@example.com"},
				{ID: "2", Name: "Hossain", Email: "hossain@example.com"},
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)

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
	expected := model.User{
		ID:    "1",
		Name:  "Nazibul",
		Email: "nazibul@example.com",
	}

	handler := setup(&mockUserStore{
		GetFn: func(ctx context.Context, id string) (model.User, error) {
			return model.User{ID: id, Name: "Nazibul", Email: "nazibul@example.com"}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	var userData model.User
	if err := json.NewDecoder(res.Body).Decode(&userData); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if userData != expected {
		t.Errorf("expected %+v, got %+v", expected, userData)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	handler := setup(&mockUserStore{
		GetFn: func(ctx context.Context, id string) (model.User, error) {
			return model.User{}, errors.New("not found")
		},
	})

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
	handler := setup(&mockUserStore{})

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
	// The mock simulates the store saying "yes user exists,
	// updated successfully" — you don't need to actually create the user first.
	handler := setup(&mockUserStore{
		UpdateFn: func(ctx context.Context, id string, user model.User) error {
			if id != "1" {
				t.Errorf("expected id 1, got %s", id)
			}
			if user.Name != "Updated Name" {
				t.Errorf("expected name Updated Name, got %s", user.Name)
			}
			return nil
		},
	})

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
	handler := setup(&mockUserStore{
		UpdateFn: func(ctx context.Context, id string, user model.User) error {
			return errors.New("User not found")
		},
	})
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
	handler := setup(&mockUserStore{})

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
	handler := setup(&mockUserStore{
		DeleteFn: func(ctx context.Context, id string) error {
			if id != "1" {
				t.Errorf("expected id 1, got %s", id)
			}
			return nil
		},
	})

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
	handler := setup(&mockUserStore{
		DeleteFn: func(ctx context.Context, id string) error {
			return errors.New("not found")
		},
	})
	req := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()

	handler.DeleteUser(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", res.StatusCode)
	}
}
