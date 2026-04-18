# Architecture Documentation

## Overview
Simple REST API following clean architecture principles with standard library routing.
```
┌─────────────────────────────────────────┐
│         HTTP Layer (main.go)            │
│    Server config + Route registration   │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│       Handler Layer (handler/)          │
│   Request parsing, Response formatting  │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│        Store Layer (store/)             │
│    In-memory storage with Mutex       │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│       Model Layer (model/)              │
│         Data structures only            │
└─────────────────────────────────────────┘
```

## Components

### 1. Server Layer (`internal/app/server.go`)
**Responsibility:** HTTP server configuration and lifecycle
```go
func NewServer(store *store.UserStore) *http.Server {
    mux := http.NewServeMux()
    handler := handler.NewUserHandler(store)
    
    // Route registration
    mux.HandleFunc("GET /users", handler.GetAllUsers)
    mux.HandleFunc("GET /users/{id}", handler.GetUser)
    mux.HandleFunc("POST /users", handler.CreateUser)
    mux.HandleFunc("PUT /users/{id}", handler.UpdateUser)
    mux.HandleFunc("DELETE /users/{id}", handler.DeleteUser)
    
    return &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }
}
```

**Key Features:**
- Production-ready timeouts
- Go 1.22+ method-based routing (`GET /users`, `POST /users`)
- Clean separation from main.go

---

### 2. Handler Layer (`internal/handler/user_handler.go`)
**Responsibility:** HTTP request/response handling
```go
type UserHandler struct {
    store *store.UserStore
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request
    // 2. Validate input
    // 3. Call store
    // 4. Format response
}
```

**Handles:**
- JSON parsing/encoding
- HTTP status codes
- Error responses
- Input validation

---

### 3. Store Layer (`internal/store/user_store.go`)
**Responsibility:** Data persistence (in-memory)
```go
type UserStore struct {
    mu    sync.Locker
    users map[string]model.User
}
```

**Thread Safety:**
- `RLock()` for reads (GetUser, GetAllUsers)
- `Lock()` for writes (CreateUser, UpdateUser, DeleteUser)
- Allows concurrent reads, exclusive writes

**Operations:**
| Operation | Lock Type | Time Complexity |
|-----------|-----------|-----------------|
| GetUser | Mutex | O(1) |
| GetAllUsers | Mutex | O(n) |
| CreateUser | Mutex | O(1) |
| UpdateUser | Mutex | O(1) |
| DeleteUser | Mutex | O(1) |

---
### Future Evolution
If the application becomes read-heavy or introduces higher concurrency,
`sync.Mutex` can be replaced with `sync.RWMutex` without changing
the public store API or handler logic.


### 4. Model Layer (`internal/model/user_model.go`)
**Responsibility:** Data structures
```go
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

**Pure data - no logic**

---

## Data Flow Example

### Creating a User
```
1. Client sends: POST /users
   Body: {"id":"1","name":"John","email":"john@example.com","age":30}
   
2. HTTP Server (ServeMux)
   ↓ Routes to handler based on "POST /users"
   
3. UserHandler.CreateUser()
   ↓ json.NewDecoder(r.Body).Decode(&user)
   ↓ Basic validation
   
4. UserStore.CreateUser()
   ↓ store.mu.Lock()
   ↓ Check duplicate ID
   ↓ store.users[user.ID] = user
   ↓ store.mu.Unlock()
   
5. Response
   ↓ w.WriteHeader(http.StatusCreated)
   ↓ json.NewEncoder(w).Encode(user)
```

---


## Concurrency Model

### Why RWMutex?

**Scenario:** 100 concurrent requests
- 90 GET requests (reads)
- 10 POST requests (writes)

**With Mutex (old approach):**
```
All 100 requests execute sequentially
Total time: ~100ms (1ms each)
```

**With RWMutex (current):**
```
90 reads execute in parallel
10 writes execute one at a time
Total time: ~11ms (1ms reads in parallel + 10ms writes)
```

**Performance gain:** ~9x faster for read-heavy workloads

---

