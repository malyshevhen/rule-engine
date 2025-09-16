package api

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		slog.Info("HTTP request",
			"method", r.Method,
			"url", r.URL.Path,
			"duration", duration,
		)
	})
}

// AuthMiddleware is a placeholder for authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement JWT or API key authentication
		// For now, allow all requests
		next.ServeHTTP(w, r)
	})
}
