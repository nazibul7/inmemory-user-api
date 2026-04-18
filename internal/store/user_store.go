package store

import (
	"context"
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

func (s *UserStore) CreateUser(ctx context.Context, user model.User) error {

	// ctx.Done() = "wait until cancelled"
	// ctx.Err()  = "check if already cancelled"
	//
	//
	// Why ctx.Err() and NOT select + ctx.Done():
	//
	// select + goroutine + channel is for RACING slow work against a timeout:
	//
	//   go func() { slowDBCall() }()     // could take 100ms-seconds
	//   select {
	//   case <-ctx.Done(): ...           // fires if DB is too slow
	//   case res := <-done: ...          // fires if DB responds in time
	//   }
	//
	// In-memory map operations take NANOSECONDS.
	// There is nothing to race — the work always finishes
	// before any timeout could ever fire.
	//
	// ctx.Err() is a simple non-blocking check:
	// "was this context already cancelled before I start?"
	// If yes  → return early, skip the work
	// If no   → proceed, work completes in nanoseconds
	//
	// select + ctx.Done() here would add:
	// → goroutine per request (unnecessary allocation)
	// → channel per request (unnecessary allocation)
	// → complexity with zero real benefit

	if err := ctx.Err(); err != nil {
		return err
	}

	// Never hold a lock to check something that doesn't require the lock that's why put ctx.Err outside of lock
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[user.ID]; exists {
		return errors.New("User already exists")
	}
	s.users[user.ID] = user
	return nil
}

func (s *UserStore) GetUser(ctx context.Context, id string) (model.User, error) {

	if err := ctx.Err(); err != nil {
		return model.User{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[id]
	if !ok {
		return model.User{}, errors.New("User not found")
	}
	return user, nil
}

func (s *UserStore) GetAllUser(ctx context.Context) ([]model.User, error) {

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	users := make([]model.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

func (s *UserStore) UpdateUser(ctx context.Context, id string, user model.User) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[id]; !exists {
		return errors.New("User not found")
	}
	s.users[id] = user
	return nil
}

func (s *UserStore) DeleteUser(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return errors.New("User not found")
	}
	delete(s.users, id)
	return nil
}
