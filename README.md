# In-Memory User API (GO)

A simple REST Api written in Go using standard library **net/http** that performs CURD operations on an **in-memory user database**.

This project demonstrates:
- Clean Go project structure
- In-memory data storage using `map` + `sync.RWMutex`
- RESTful CRUD APIs
- Go 1.22+ `ServeMux` method-based routing
- Production-style server configuration (timeouts, explicit wiring)

---
## Project Structure
```
inmemory-user-api/
├── go.mod
├── main.go
├── internal/
│ ├── app/
│ │ └── server.go # HTTP server setup
│ ├── handler/
│ │ └── user_handler.go # HTTP handlers (CRUD)
│ ├── model/
│ │ └── user_model.go # User model
│ └── store/
│ └── user_store.go # In-memory user store
└── README.md
```
---

## Running the Server
```
go run main.go

```
