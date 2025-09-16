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
	"github.com/malyshevhen/rule-engine/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
)

var (
	// Simple rate limiter: 100 requests per minute per IP
	// Using sync.Map for better concurrency
	requestCounts sync.Map // map[string]int
	lastResets    sync.Map // map[string]time.Time
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

// TracingMiddleware adds distributed tracing to HTTP requests
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracing.StartSpan(r.Context(), "http.request")
		defer span.End()

		// Add HTTP attributes to the span
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.Path),
			attribute.String("http.user_agent", r.Header.Get("User-Agent")),
		)

		// Add span to request context
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
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
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
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
		if after, ok := strings.CutPrefix(authHeader, "ApiKey "); ok {
			apiKey := after
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

		now := time.Now()

		// Get or initialize rate limit data for this IP
		lastResetVal, _ := lastResets.Load(ip)
		var lastReset time.Time
		if lastResetVal != nil {
			lastReset = lastResetVal.(time.Time)
		}

		countVal, _ := requestCounts.Load(ip)
		var count int
		if countVal != nil {
			count = countVal.(int)
		}

		// Reset counter if window has passed
		if lastReset.IsZero() || now.Sub(lastReset) >= window {
			count = 0
			lastResets.Store(ip, now)
		}

		// Check rate limit
		if count >= maxRequests {
			ErrorResponse(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}

		// Increment counter
		requestCounts.Store(ip, count+1)

		next.ServeHTTP(w, r)
	})
}

// DisableRateLimiting disables rate limiting (for performance testing)
func DisableRateLimiting() {
	rateLimitingEnabled = false
}

// EnableRateLimiting enables rate limiting
func EnableRateLimiting() {
	rateLimitingEnabled = true
}

// AuthMiddleware supports both JWT and API key authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return APIKeyMiddleware(next)
}
