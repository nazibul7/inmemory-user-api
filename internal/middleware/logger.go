package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := r.Context().Value(RequestIDKey).(string)
		next.ServeHTTP(w, r)

		log.Printf("requestID=%v methode=%s URL=%s duration=%s", reqID, r.Method, r.URL, time.Since(start))
	})
}
