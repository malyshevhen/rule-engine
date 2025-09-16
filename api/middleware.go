package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/malyshevhen/rule-engine/internal/storage/redis"
	"github.com/malyshevhen/rule-engine/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
)

var (
	// Redis rate limiter: 100 requests per minute per IP
	redisRateLimiter *redis_rate.Limiter
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

// RateLimitMiddleware limits requests per IP (100 per minute) using Redis-backed rate limiting
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rateLimitingEnabled || redisRateLimiter == nil {
			next.ServeHTTP(w, r)
			return
		}

		ip := r.RemoteAddr
		// Simple IP extraction, in production use proper IP extraction

		// Check rate limit using Redis
		result, err := redisRateLimiter.Allow(r.Context(), ip, redis_rate.PerMinute(100))
		if err != nil {
			slog.Error("Rate limiter error", "error", err)
			// Allow request on error to avoid blocking legitimate traffic
			next.ServeHTTP(w, r)
			return
		}

		// Check if request is allowed
		if result.Allowed == 0 {
			ErrorResponse(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// InitRedisRateLimiter initializes the Redis-backed rate limiter
func InitRedisRateLimiter(redisClient *redis.Client) {
	redisRateLimiter = redis_rate.NewLimiter(redisClient.GetClient())
}

// GetRedisRateLimiter returns the Redis rate limiter instance
func GetRedisRateLimiter() *redis_rate.Limiter {
	return redisRateLimiter
}

// ResetMiddlewareForTesting resets the middleware state for testing
// This should only be used in test code to ensure test isolation
func ResetMiddlewareForTesting() {
	// Reset rate limiting enabled flag
	rateLimitingEnabled = true
	// Note: Redis rate limiter state is managed by Redis itself and persists across tests
	// In a real testing scenario, you might want to flush Redis or use a test-specific Redis instance
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
