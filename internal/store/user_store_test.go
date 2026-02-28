package store

import (
	"testing"

	"github.com/nazibul7/inmemory-user-api/internal/model"
)

func newUserStore() *UserStore {
	return NewUserStore()
}

func newTestUser(id string) model.User {
	return model.User{
		ID:    id,
		Name:  "Test User",
		Email: "test-user@gmail.com",
	}
}

// -----------CreateUser-------------------------------------------------
func TestCreateUser_Success(t *testing.T) {
	s := newUserStore()
	user := newTestUser("1")

	if err := s.CreateUser(user); err != nil {
		t.Fatalf("Expected no error,got %v", err)
	}
}

func TestCreateUser_Duplicate(t *testing.T) {
	s := newUserStore()
	user := newTestUser("1")
	s.CreateUser(user)

	if err := s.CreateUser(user); err == nil { // Using same ID again
		t.Fatal("expected error for duplicate user, got nil")
	}
}

// -----------GetUser-------------------------------------------------
func TestGetUser_Success(t *testing.T) {
	s := newUserStore()
	user := newTestUser("1")
	s.CreateUser(user)

	got, err := s.GetUser("1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.ID != user.ID {
		t.Errorf("expected ID %s, got %s", user.ID, got.ID)
	}
}

func TestUser_NotFound(t *testing.T) {
	s := newUserStore()

	_, err := s.GetUser("doesnotexist")
	if err == nil {
		t.Fatal("expected error for missing user, got nil")
	}
}

// -----------GetUAlluser-------------------------------------------------
func TestGetAllUser_Empty(t *testing.T) {
	s := newUserStore()

	users := s.GetAllUser()
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

func TestGetAllUser_ReturnsAll(t *testing.T) {
	s := newUserStore()
	s.CreateUser(newTestUser("1"))
	s.CreateUser(newTestUser("2"))
	s.CreateUser(newTestUser("3"))

	users := s.GetAllUser()
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}

//---------UpdateUser--------------------------------------------

func TestUpdateUser_Success(t *testing.T) {
	s := newUserStore()
	s.CreateUser(newTestUser("1"))

	updated := model.User{ID: "1", Name: "Updated Name", Email: "updated@example.com"}
	if err := s.UpdateUser("1", updated); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, err := s.GetUser("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != updated.Name {
		t.Errorf("expected updated name, got %s", got.Name)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	s := newUserStore()

	err := s.UpdateUser("ghost", newTestUser("ghost"))
	if err == nil {
		t.Fatal("expected error for missing user, got nil")
	}
}

// -----DeleteUser-----------------------------------------------
func TestDeleteUser_Success(t *testing.T) {
	s := newUserStore()
	s.CreateUser(newTestUser("1"))

	err := s.DeleteUser("1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = s.GetUser("1")
	if err == nil {
		t.Fatal("expected user to be deleted, but still found")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	s := newUserStore()

	err := s.DeleteUser("ghost")
	if err == nil {
		t.Fatal("expected error for missing user, got nil")
	}
}
