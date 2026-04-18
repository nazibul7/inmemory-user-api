## Design Decisions

### 1. Standard Library Routing (Go 1.22+)

**Why not Gorilla Mux or Chi?**
- ✅ Zero dependencies
- ✅ Method-based routing built-in
- ✅ Path parameters support (`{id}`)
- ✅ Fast (no reflection)
- ✅ Maintained by Go team

**Before Go 1.22:**
```go
// Needed external router
r.HandleFunc("/users", handler.GetUsers).Methods("GET")
r.HandleFunc("/users", handler.CreateUser).Methods("POST")
```

**After Go 1.22:**
```go
// Built-in method routing
mux.HandleFunc("GET /users", handler.GetUsers)
mux.HandleFunc("POST /users", handler.CreateUser)
```

---

### 2. String IDs vs Auto-increment

**Current:** `id string`

**Rationale:**
- Simple for in-memory storage
- Client controls ID generation
- Easy to test
- No UUID library needed

**When adding database:**
- Consider UUID v4 server-side generation
- Or switch to auto-increment integers

---

### 3. No Service Layer

**Current structure:**
```
Handler → Store
```

**When to add Service layer:**
```
Handler → Service → Store
```

**Add Service layer when you need:**
- Multi-step business logic
- Transaction coordination
- Cross-entity operations
- External API calls
- Complex validation rules

**Current app:** Simple CRUD = no service layer needed

---

## Error Handling Strategy

### Current Approach
```go
// Store errors
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
)

// Handler converts to HTTP
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.store.GetUser(id)
    if err != nil {
        if errors.Is(err, store.ErrUserNotFound) {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }
    // Success path
}
```

**Mapping:**
| Store Error | HTTP Status | Response |
|-------------|-------------|----------|
| ErrUserNotFound | 404 | "user not found" |
| ErrUserAlreadyExists | 409 | "user already exists" |
| Any other error | 500 | "Internal error" |

---

## Server Configuration

### Production-Ready Timeouts
```go
&http.Server{
    ReadTimeout:  10 * time.Second,  // Reading request body
    WriteTimeout: 10 * time.Second,  // Writing response
    IdleTimeout:  120 * time.Second, // Keep-alive connections
}
```

**Why these matter:**
- Prevents slow client attacks
- Manages connection pool size
- Graceful resource cleanup

**Without timeouts:**
- Slow clients can exhaust connections
- Memory leaks from abandoned connections
- Vulnerable to Slowloris attacks

---

## Future Evolution Path

### Phase 1: Testing (Next step)
```
Add:
- Unit tests for store
- Handler tests with httptest
- Integration tests
```

### Phase 2: Database
```
Add:
- PostgreSQL store implementation
- Interface abstraction
- Migration scripts
```

### Phase 3: Production Features
```
Add:
- Structured logging (slog)
- Metrics (Prometheus)
- Health check endpoint
- Graceful shutdown
- Docker containerization
```

### Phase 4: Advanced Features
```
Add:
- Authentication (JWT)
- Input validation (validator)
- Rate limiting
- API versioning
- OpenAPI docs
```

---

## Performance Characteristics

### Current Benchmarks (estimated)

| Metric | Value | Notes |
|--------|-------|-------|
| Requests/sec | ~10,000 | Single instance |
| Latency (p50) | <1ms | In-memory |
| Latency (p99) | <5ms | Lock contention |
| Memory usage | ~1KB per user | Go map overhead |
| Max users | ~1M | Limited by RAM |

**Bottleneck:** Single instance, no horizontal scaling

---

## Code Quality Standards

### Current State
✅ Clean separation of concerns
✅ Thread-safe concurrent access
✅ Production-ready server config
✅ Standard library only (no deps)
✅ Go 1.22+ modern routing

### To Add
⬜ Unit tests (>80% coverage)
⬜ Integration tests
⬜ Structured logging
⬜ Input validation
⬜ Error wrapping with context

---

## References

- [Go 1.22 Routing](https://go.dev/blog/routing-enhancements)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go HTTP Server Guide](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)
EOF