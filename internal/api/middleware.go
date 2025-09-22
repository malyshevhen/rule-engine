package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
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

// recoveryMiddleware recovers from panics and returns a 500 error
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Panic recovered in HTTP handler",
					"panic", err,
					"method", r.Method,
					"url", r.URL.Path,
					"stack", string(debug.Stack()),
				)
				ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
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

// tracingMiddleware adds distributed tracing to HTTP requests
func tracingMiddleware(next http.Handler) http.Handler {
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

// jwtMiddleware validates JWT tokens for authentication
// TODO: Implement JWT authentication middleware when needed
// nolint:unused
func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			ErrorResponse(w, http.StatusUnauthorized, "AUTHENTICATION_ERROR", "Missing authorization header")
			return
		}

		// Extract token from "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ErrorResponse(w, http.StatusUnauthorized, "AUTHENTICATION_ERROR", "Invalid authorization header format")
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
			ErrorResponse(w, http.StatusUnauthorized, "AUTHENTICATION_ERROR", "Invalid token")
			return
		}

		if !token.Valid {
			ErrorResponse(w, http.StatusUnauthorized, "AUTHENTICATION_ERROR", "Invalid token")
			return
		}

		// Token is valid, proceed
		next.ServeHTTP(w, r)
	})
}

// authenticateRequest validates both API key and JWT authentication
// Either API key (X-API-Key header) or JWT Bearer token (Authorization header) must be valid
func authenticateRequest(r *http.Request) error {
	// Check API key authentication first
	apiKey := r.Header.Get("X-API-Key")
	if apiKey != "" {
		expectedKey := os.Getenv("API_KEY")
		if expectedKey == "" {
			slog.Error("API_KEY environment variable not set")
			return fmt.Errorf("authentication configuration error")
		}
		if apiKey == expectedKey {
			return nil // Valid API key
		}
	}

	// Check JWT Bearer token authentication
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Extract token from "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// Not a Bearer token format, continue to check if we have a valid API key
			if apiKey != "" {
				return fmt.Errorf("invalid API key")
			}
			return fmt.Errorf("missing valid authentication")
		}

		// Parse and validate JWT token
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
			return fmt.Errorf("invalid token")
		}

		if token.Valid {
			return nil // Valid JWT token
		}
	}

	// Neither API key nor JWT token was valid
	return fmt.Errorf("missing or invalid authentication")
}

// rateLimitMiddleware limits requests per IP (100 per minute) using Redis-backed rate limiting
func rateLimitMiddleware(next http.Handler) http.Handler {
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
			ErrorResponse(w, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Rate limit exceeded")
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

// AuthMiddleware supports both JWT Bearer token and API key authentication
// Either X-API-Key header or Authorization: Bearer <token> header must be valid
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := authenticateRequest(r); err != nil {
			ErrorResponse(w, http.StatusUnauthorized, "AUTHENTICATION_ERROR", err.Error())
			return
		}

		// Authentication successful, proceed
		next.ServeHTTP(w, r)
	})
}
