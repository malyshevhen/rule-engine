package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// Simple rate limiter: 100 requests per minute per IP
	requestCounts = make(map[string]int)
	lastResets    = make(map[string]time.Time)
	rateLimitMu   sync.Mutex
	maxRequests   = 100
	window        = time.Minute
	// Allow disabling rate limiting for performance tests
	rateLimitingEnabled = true
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

// JWTMiddleware validates JWT tokens for authentication
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			ErrorResponse(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		// Extract token from "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ErrorResponse(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				return nil, fmt.Errorf("JWT_SECRET not set")
			}
			return []byte(secret), nil
		})

		if err != nil {
			slog.Error("JWT parsing error", "error", err)
			ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		if !token.Valid {
			ErrorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Token is valid, proceed
		next.ServeHTTP(w, r)
	})
}

// APIKeyMiddleware validates API key authentication
func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			ErrorResponse(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		// Check for API key format: "ApiKey <key>"
		if strings.HasPrefix(authHeader, "ApiKey ") {
			apiKey := strings.TrimPrefix(authHeader, "ApiKey ")
			expectedKey := os.Getenv("API_KEY")
			if expectedKey == "" {
				slog.Error("API_KEY environment variable not set")
				ErrorResponse(w, http.StatusInternalServerError, "Authentication configuration error")
				return
			}
			if apiKey != expectedKey {
				ErrorResponse(w, http.StatusUnauthorized, "Invalid API key")
				return
			}
			// Valid API key, proceed
			next.ServeHTTP(w, r)
			return
		}

		// If not API key, try JWT
		JWTMiddleware(next).ServeHTTP(w, r)
	})
}

// RateLimitMiddleware limits requests per IP (100 per minute)
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rateLimitingEnabled {
			next.ServeHTTP(w, r)
			return
		}

		ip := r.RemoteAddr
		// Simple IP extraction, in production use proper IP extraction

		rateLimitMu.Lock()
		now := time.Now()
		lastReset, exists := lastResets[ip]
		if !exists || now.Sub(lastReset) >= window {
			requestCounts[ip] = 0
			lastResets[ip] = now
		}

		count := requestCounts[ip]
		if count >= maxRequests {
			rateLimitMu.Unlock()
			ErrorResponse(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}

		requestCounts[ip] = count + 1
		rateLimitMu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// DisableRateLimiting disables rate limiting (for performance testing)
func DisableRateLimiting() {
	rateLimitMu.Lock()
	defer rateLimitMu.Unlock()
	rateLimitingEnabled = false
}

// EnableRateLimiting enables rate limiting
func EnableRateLimiting() {
	rateLimitMu.Lock()
	defer rateLimitMu.Unlock()
	rateLimitingEnabled = true
}

// AuthMiddleware supports both JWT and API key authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return APIKeyMiddleware(next)
}
