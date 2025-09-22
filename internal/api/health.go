package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	redisClient "github.com/malyshevhen/rule-engine/internal/storage/redis"
)

// Health represents the health check service
type Health struct {
	db    *pgxpool.Pool
	redis *redisClient.Client
}

// NewHealth creates a new Health instance
func NewHealth(db *pgxpool.Pool, redis *redisClient.Client) *Health {
	return &Health{db: db, redis: redis}
}

// @Summary Health check
// @Description Check the health status of the service and its dependencies.
// @Tags health
// @Produce  json
// @Success 200  {object}  map[string]string
// @Failure 503  {object}  map[string]string
// @Router /health [get]
// healthCheckHandler returns a handler for the health check endpoint
func (h *Health) healthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		dbStatus := "ok"
		if err := h.db.Ping(ctx); err != nil {
			dbStatus = "error"
		}

		redisStatus := "ok"
		if h.redis != nil {
			if err := h.redis.Ping(ctx); err != nil {
				redisStatus = "error"
			}
		} else {
			redisStatus = "disabled"
		}

		response := map[string]string{
			"database": dbStatus,
			"redis":    redisStatus,
		}

		w.Header().Set("Content-Type", "application/json")
		if dbStatus != "ok" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
