package main

import (
	"log"
	"net/http"

	"github.com/nazibul7/inmemory-user-api/internal/app"
	"github.com/nazibul7/inmemory-user-api/internal/handler"
	"github.com/nazibul7/inmemory-user-api/internal/store"
)

func main() {
	UserStore := store.NewUserStore()
	UserHandler := handler.NewUserHandler(UserStore)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", UserHandler.GetAll)
	mux.HandleFunc("POST /user", UserHandler.Create)
	mux.HandleFunc("GET /user/{id}", UserHandler.GetUser)
	mux.HandleFunc("PUT /user/{id}", UserHandler.UpdateUser)
	mux.HandleFunc("DELETE /user/{id}", UserHandler.DeleteUser)


	server := app.NewServer(":9000", mux)

	log.Println("Server running on :9000")
	log.Fatal(server.ListenAndServe())
}
