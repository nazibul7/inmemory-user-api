package store

import (
	"context"
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

	ctx := context.Background()
	if err := s.CreateUser(ctx, user); err != nil {
		t.Fatalf("Expected no error,got %v", err)
	}
}

func TestCreateUser_Duplicate(t *testing.T) {
	s := newUserStore()
	user := newTestUser("1")
	ctx := context.Background()
	s.CreateUser(ctx, user)

	if err := s.CreateUser(ctx, user); err == nil { // Using same ID again
		t.Fatal("expected error for duplicate user, got nil")
	}
}

// -----------GetUser-------------------------------------------------
func TestGetUser_Success(t *testing.T) {
	s := newUserStore()
	user := newTestUser("1")
	ctx := context.Background()
	s.CreateUser(ctx, user)

	got, err := s.GetUser(ctx, "1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.ID != user.ID {
		t.Errorf("expected ID %s, got %s", user.ID, got.ID)
	}
}

func TestUser_NotFound(t *testing.T) {
	s := newUserStore()

	ctx := context.Background()
	_, err := s.GetUser(ctx, "doesnotexist")
	if err == nil {
		t.Fatal("expected error for missing user, got nil")
	}
}

// -----------GetUAlluser-------------------------------------------------
func TestGetAllUser_Empty(t *testing.T) {
	s := newUserStore()

	ctx := context.Background()
	users, err := s.GetAllUser(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

func TestGetAllUser_ReturnsAll(t *testing.T) {
	s := newUserStore()
	ctx := context.Background()
	s.CreateUser(ctx, newTestUser("1"))
	s.CreateUser(ctx, newTestUser("2"))
	s.CreateUser(ctx, newTestUser("3"))

	users, err := s.GetAllUser(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}

//---------UpdateUser--------------------------------------------

func TestUpdateUser_Success(t *testing.T) {
	s := newUserStore()
	ctx := context.Background()
	s.CreateUser(ctx, newTestUser("1"))

	updated := model.User{ID: "1", Name: "Updated Name", Email: "updated@example.com"}
	if err := s.UpdateUser(ctx, "1", updated); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, err := s.GetUser(ctx, "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != updated.Name {
		t.Errorf("expected updated name, got %s", got.Name)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	s := newUserStore()

	ctx := context.Background()
	err := s.UpdateUser(ctx, "ghost", newTestUser("ghost"))
	if err == nil {
		t.Fatal("expected error for missing user, got nil")
	}
}

// -----DeleteUser-----------------------------------------------
func TestDeleteUser_Success(t *testing.T) {
	s := newUserStore()

	ctx := context.Background()
	s.CreateUser(ctx, newTestUser("1"))

	err := s.DeleteUser(ctx, "1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = s.GetUser(ctx, "1")
	if err == nil {
		t.Fatal("expected user to be deleted, but still found")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	s := newUserStore()

	ctx := context.Background()
	err := s.DeleteUser(ctx, "ghost")
	if err == nil {
		t.Fatal("expected error for missing user, got nil")
	}
}
