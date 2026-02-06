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

/** Here sync.Locker is an interface.
   type Locker interface {
	 Lock()
	 Unlock()
   }

   sync.Mutex is
   type Mutex struct {
	state int32
	sema  uint32
   }
   it has methods
   func (m *Mutex) Lock()
   func (m *Mutex) Unlock()

   That's why used &sync.Mutex{} to satisfy interface.

   Now question is if it's implements value receiver it will not work because
   if the methods used value receivers, Lock would get a copy and it would not sync properly.

   *** Pointer to Interface is Almost Always Useless-
   type Locker interface {
    Lock()
    Unlock()
}

 ❌ USELESS - Pointer to interface
 var mu *Locker

 ✅ CORRECT - Interface variable (not pointer to interface)
 var mu Locker

 Interfaces Already Hold Pointers Internally!
 An interface value consists of two pointers:
 go// Internal representation of an interface
 type interface {
    type  *_type      // Pointer to type information
    data  unsafe.Pointer  // Pointer to the actual value
 }
*/

func NewUserStore() *UserStore {
	return &UserStore{
		mu:    &sync.Mutex{},
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
	if _, exists := s.users[id]; !exists {
		return errors.New("User not found")
	}
	s.users[id] = user
	return nil
}

func (s *UserStore) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return errors.New("User not found")
	}
	delete(s.users, id)
	return nil
}
