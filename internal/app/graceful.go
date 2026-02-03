package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunWithGracefulShutdown(server *http.Server, timeout time.Duration) error {
	// Channel to listen for errors from listeners
	serverError := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on %s", server.Addr)
		serverError <- server.ListenAndServe()
	}()

	// Channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(shutdown)

	// Block until we receive an error or signal
	select {
	case err := <-serverError:
		log.Printf("Error starting server: %v", err)
		return err
	case sig := <-shutdown:
		log.Printf("Shutdown signal received: %v", sig)
		log.Println("Starting graceful shutdown...")

		// Create context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Attempt graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			log.Println("Forcing server to close...")

			// Force close if graceful shutdown fails
			if closeErr := server.Close(); closeErr != nil {
				// Return both errors
				return errors.Join(err, closeErr)
			}
			return err
		}
		log.Println("Server stopped gracefully")
	}
	return nil
}
