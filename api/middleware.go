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
	"golang.org/x/time/rate"
)

var (
	// Rate limiter: 100 requests per minute per IP using token bucket algorithm
	limiters   = make(map[string]*limiterEntry)
	limitersMu sync.RWMutex
	rateLimit  = rate.Limit(100.0 / 60.0) // 100 requests per minute
	burstLimit = 100                      // Allow bursts of up to 100 requests for testing
	// Allow disabling rate limiting for performance tests
	rateLimitingEnabled = true
	// Control cleanup goroutine
	cleanupEnabled = false
	cleanupDone    = make(chan struct{})
)

type limiterEntry struct {
	limiter  *rate.Limiter
	lastUsed time.Time
}

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

// RateLimitMiddleware limits requests per IP (100 per minute) using token bucket algorithm
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rateLimitingEnabled {
			next.ServeHTTP(w, r)
			return
		}

		ip := r.RemoteAddr
		// Simple IP extraction, in production use proper IP extraction

		// Get or create rate limiter for this IP
		limiter := getLimiter(ip)

		// Check if request is allowed
		if !limiter.Allow() {
			ErrorResponse(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getLimiter returns the rate limiter for the given IP, creating one if it doesn't exist
func getLimiter(ip string) *rate.Limiter {
	limitersMu.Lock()
	defer limitersMu.Unlock()

	entry, exists := limiters[ip]
	if !exists {
		entry = &limiterEntry{
			limiter:  rate.NewLimiter(rateLimit, burstLimit),
			lastUsed: time.Now(),
		}
		limiters[ip] = entry
	} else {
		entry.lastUsed = time.Now()
	}

	return entry.limiter
}

// cleanupLimiters removes rate limiters that haven't been used in the last hour
func cleanupLimiters() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			limitersMu.Lock()
			for ip, entry := range limiters {
				if time.Since(entry.lastUsed) > 1*time.Hour {
					delete(limiters, ip)
				}
			}
			limitersMu.Unlock()
		case <-cleanupDone:
			return
		}
	}
}

// StartRateLimiterCleanup starts the background cleanup of unused rate limiters
func StartRateLimiterCleanup() {
	if !cleanupEnabled {
		cleanupEnabled = true
		go cleanupLimiters()
	}
}

// StopRateLimiterCleanup stops the background cleanup of unused rate limiters
func StopRateLimiterCleanup() {
	if cleanupEnabled {
		cleanupEnabled = false
		close(cleanupDone)
	}
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
