package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"runtime"

	"github.com/nazibul7/inmemory-user-api/internal/app"
	"github.com/nazibul7/inmemory-user-api/internal/handler"
	"github.com/nazibul7/inmemory-user-api/internal/middleware"
	"github.com/nazibul7/inmemory-user-api/internal/store"
)

func main() {
	UserStore := store.NewUserStore()
	UserHandler := handler.NewUserHandler(UserStore)

	mux := http.NewServeMux()
	muxHandler := middleware.RequestID(mux)
	muxHandler = middleware.Logger(muxHandler)
	muxHandler=middleware.Recoverer(muxHandler)

	mux.HandleFunc("GET /users", UserHandler.GetAll)
	mux.HandleFunc("POST /user", UserHandler.Create)
	mux.HandleFunc("GET /user/{id}", UserHandler.GetUser)
	mux.HandleFunc("PUT /user/{id}", UserHandler.UpdateUser)
	mux.HandleFunc("DELETE /user/{id}", UserHandler.DeleteUser)

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	server := app.NewServer(":9000", muxHandler)

	if err := app.RunWithGracefulShutdown(server, 30); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	if runtime.NumGoroutine() > 1 {
		fmt.Println("\n Leaked goroutine stack trace:")
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, true)
		fmt.Printf("%s\n", buf[:stackSize])
	}
}
