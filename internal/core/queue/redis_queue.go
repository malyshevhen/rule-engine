package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	redisClient "github.com/malyshevhen/rule-engine/internal/storage/redis"
	"github.com/redis/go-redis/v9"
)

// RedisQueue implements a persistent queue using Redis
type RedisQueue struct {
	client     *redisClient.Client
	queueKey   string
	processing bool
}

// NewRedisQueue creates a new Redis-backed queue
func NewRedisQueue(client *redisClient.Client, queueKey string) *RedisQueue {
	return &RedisQueue{
		client:   client,
		queueKey: queueKey,
	}
}

// Enqueue adds a rule execution request to the Redis queue
func (q *RedisQueue) Enqueue(ctx context.Context, req *ExecutionRequest) error {
	if q.client == nil {
		return fmt.Errorf("Redis client not available")
	}

	req.ID = uuid.New()
	req.QueuedAt = time.Now()

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Use score as timestamp for ordering (FIFO)
	score := float64(req.QueuedAt.UnixNano())

	err = q.client.Set(ctx, fmt.Sprintf("queue:item:%s", req.ID.String()), string(data), 24*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to store request: %w", err)
	}

	// Add to sorted set for ordering
	err = q.client.GetClient().ZAdd(ctx, q.queueKey, redis.Z{
		Score:  score,
		Member: req.ID.String(),
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to add to queue: %w", err)
	}

	slog.Debug("Enqueued rule execution request", "request_id", req.ID, "rule_id", req.RuleID)
	return nil
}

// Dequeue removes and returns the next rule execution request from the Redis queue
func (q *RedisQueue) Dequeue(ctx context.Context) (*ExecutionRequest, error) {
	if q.client == nil {
		return nil, fmt.Errorf("Redis client not available")
	}

	// Get the first item from the sorted set (lowest score = oldest)
	result := q.client.GetClient().ZPopMin(ctx, q.queueKey, 1)
	if result.Err() != nil {
		return nil, result.Err()
	}

	vals := result.Val()
	if len(vals) == 0 {
		return nil, ErrQueueEmpty
	}

	requestID := vals[0].Member.(string)

	// Get the actual request data
	data, err := q.client.Get(ctx, fmt.Sprintf("queue:item:%s", requestID))
	if err != nil {
		// If we can't get the data, log error but continue (item might have expired)
		slog.Error("Failed to get queued request data", "request_id", requestID, "error", err)
		return nil, ErrQueueEmpty
	}

	var req ExecutionRequest
	if err := json.Unmarshal([]byte(data), &req); err != nil {
		slog.Error("Failed to unmarshal queued request", "request_id", requestID, "error", err)
		// Clean up corrupted data
		q.client.Del(ctx, fmt.Sprintf("queue:item:%s", requestID))
		return nil, ErrQueueEmpty
	}

	// Clean up the stored data
	q.client.Del(ctx, fmt.Sprintf("queue:item:%s", requestID))

	slog.Debug("Dequeued rule execution request", "request_id", req.ID, "rule_id", req.RuleID)
	return &req, nil
}

// Size returns the current number of items in the Redis queue
func (q *RedisQueue) Size() int {
	if q.client == nil {
		return 0
	}

	ctx := context.Background()
	size, err := q.client.GetClient().ZCard(ctx, q.queueKey).Result()
	if err != nil {
		slog.Error("Failed to get queue size", "error", err)
		return 0
	}

	return int(size)
}

// Close closes the Redis queue (no-op for Redis implementation)
func (q *RedisQueue) Close() error {
	// Redis connections are managed externally
	return nil
}

// CleanupExpired removes expired queue items (older than specified duration)
func (q *RedisQueue) CleanupExpired(ctx context.Context, maxAge time.Duration) error {
	if q.client == nil {
		return nil
	}

	cutoff := time.Now().Add(-maxAge).UnixNano()

	// Remove items older than cutoff from the sorted set
	removed, err := q.client.GetClient().ZRemRangeByScore(ctx, q.queueKey, "-inf", strconv.FormatInt(cutoff, 10)).Result()
	if err != nil {
		return fmt.Errorf("failed to cleanup expired items: %w", err)
	}

	if removed > 0 {
		slog.Info("Cleaned up expired queue items", "count", removed)
	}

	return nil
}
