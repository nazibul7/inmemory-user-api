package store

import (
	"errors"
	"sync"

	"github.com/nazibul7/inmemory-user-api/internal/model"
)

type UserStore struct {
	mu    sync.Locker
	users map[string]model.User
}

func NewUserStore() *UserStore {
	return &UserStore{
		mu: &sync.Mutex{},
		users: make(map[string]model.User),
	}
}

func (s *UserStore) CreateUser(user model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[user.ID]; exists {
		return errors.New("User already exists")
	}
	s.users[user.ID] = user
	return nil
}

func (s *UserStore) GetUser(id string) (model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[id]
	if !ok {
		return model.User{}, errors.New("User not found")
	}
	return user, nil
}

func (s *UserStore) GetAllUser() []model.User {
	s.mu.Lock()
	defer s.mu.Unlock()

	users := make([]model.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

func (s *UserStore) UpdateUser(id string, user model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[id]; !ok {
		return errors.New("User not found")
	}
	s.users[id] = user
	return nil
}

func (s *UserStore) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[id]; !ok {
		return errors.New("User not found")
	}
	delete(s.users, id)
	return nil
}
