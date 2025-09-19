package queue

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	redisClient "github.com/malyshevhen/rule-engine/internal/storage/redis"
	"github.com/redis/go-redis/v9"
)

// generateInstanceID creates a unique identifier for this application instance
func generateInstanceID() string {
	// Generate a random component for uniqueness
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback to UUID if random generation fails
		return uuid.New().String()
	}

	// Include hostname and PID for additional uniqueness
	hostname, _ := os.Hostname()
	pid := os.Getpid()

	return fmt.Sprintf("%s-%d-%s", hostname, pid, hex.EncodeToString(randomBytes))
}

// RedisQueue implements a persistent queue using Redis
type RedisQueue struct {
	client        *redisClient.Client
	queueKey      string
	instanceID    string
	processing    bool
	metricsKey    string
	healthKey     string
	lastHeartbeat time.Time
}

// NewRedisQueue creates a new Redis-backed queue
func NewRedisQueue(client *redisClient.Client, queueKey string) *RedisQueue {
	instanceID := generateInstanceID()
	return &RedisQueue{
		client:        client,
		queueKey:      queueKey,
		instanceID:    instanceID,
		metricsKey:    queueKey + ":metrics",
		healthKey:     fmt.Sprintf("instances:%s", instanceID),
		lastHeartbeat: time.Now(),
	}
}

// Enqueue adds a rule execution request to the Redis queue
func (q *RedisQueue) Enqueue(ctx context.Context, req *ExecutionRequest) error {
	if q.client == nil {
		return fmt.Errorf("redis client not available")
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

// Dequeue removes and returns the next rule execution request from the Redis queue with distributed locking
func (q *RedisQueue) Dequeue(ctx context.Context) (*ExecutionRequest, error) {
	if q.client == nil {
		return nil, fmt.Errorf("redis client not available")
	}

	// Try to find an available item that we can lock
	maxAttempts := 10 // Limit attempts to avoid infinite loops
	for range maxAttempts {
		// Peek at the first item without removing it
		result := q.client.GetClient().ZRangeWithScores(ctx, q.queueKey, 0, 0)
		if result.Err() != nil {
			return nil, result.Err()
		}

		vals := result.Val()
		if len(vals) == 0 {
			return nil, ErrQueueEmpty
		}

		requestID := vals[0].Member.(string)

		// Try to acquire a lock for this request (30 second TTL)
		lockAcquired, err := q.acquireLock(ctx, requestID, 30*time.Second)
		if err != nil {
			slog.Error("Failed to acquire lock", "request_id", requestID, "error", err)
			continue
		}

		if !lockAcquired {
			// Another instance got the lock, try the next item
			continue
		}

		// We acquired the lock, now try to dequeue the item
		// Use ZREM to remove the specific item (in case another instance removed it)
		removed, err := q.client.GetClient().ZRem(ctx, q.queueKey, requestID).Result()
		if err != nil {
			if relErr := q.releaseLock(ctx, requestID); relErr != nil {
				slog.Warn("Failed to release lock on error", "request_id", requestID, "error", relErr)
			}
			return nil, fmt.Errorf("failed to remove item from queue: %w", err)
		}

		if removed == 0 {
			// Item was already processed by another instance
			if err := q.releaseLock(ctx, requestID); err != nil {
				slog.Warn("Failed to release lock for already processed item", "request_id", requestID, "error", err)
			}
			continue
		}

		// Get the actual request data
		data, err := q.client.Get(ctx, fmt.Sprintf("queue:item:%s", requestID))
		if err != nil {
			if relErr := q.releaseLock(ctx, requestID); relErr != nil {
				slog.Warn("Failed to release lock on data get error", "request_id", requestID, "error", relErr)
			}
			// If we can't get the data, log error but continue (item might have expired)
			slog.Error("Failed to get queued request data", "request_id", requestID, "error", err)
			continue
		}

		var req ExecutionRequest
		if err := json.Unmarshal([]byte(data), &req); err != nil {
			if relErr := q.releaseLock(ctx, requestID); relErr != nil {
				slog.Warn("Failed to release lock on unmarshal error", "request_id", requestID, "error", relErr)
			}
			slog.Error("Failed to unmarshal queued request", "request_id", requestID, "error", err)
			// Clean up corrupted data
			if delErr := q.client.Del(ctx, fmt.Sprintf("queue:item:%s", requestID)); delErr != nil {
				slog.Warn("Failed to delete corrupted queue item", "request_id", requestID, "error", delErr)
			}
			continue
		}

		// Clean up the stored data
		if err := q.client.Del(ctx, fmt.Sprintf("queue:item:%s", requestID)); err != nil {
			slog.Warn("Failed to delete processed queue item", "request_id", requestID, "error", err)
		}

		slog.Debug("Dequeued rule execution request", "request_id", req.ID, "rule_id", req.RuleID, "instance_id", q.instanceID)
		return &req, nil
	}

	// If we couldn't find an available item after max attempts, return empty
	return nil, ErrQueueEmpty
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

// acquireLock attempts to acquire a distributed lock for processing a specific request
func (q *RedisQueue) acquireLock(ctx context.Context, requestID string, lockTTL time.Duration) (bool, error) {
	if q.client == nil {
		return false, fmt.Errorf("redis client not available")
	}

	lockKey := fmt.Sprintf("lock:%s", requestID)

	// Try to set the lock with NX (only if not exists) and EX (expire)
	success, err := q.client.GetClient().SetNX(ctx, lockKey, q.instanceID, lockTTL).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock for request %s: %w", requestID, err)
	}

	if success {
		slog.Debug("Acquired distributed lock", "request_id", requestID, "instance_id", q.instanceID)
	}

	return success, nil
}

// ReleaseLock releases the distributed lock for a specific request (public method)
func (q *RedisQueue) ReleaseLock(ctx context.Context, requestID string) error {
	return q.releaseLock(ctx, requestID)
}

// releaseLock releases the distributed lock for a specific request
func (q *RedisQueue) releaseLock(ctx context.Context, requestID string) error {
	if q.client == nil {
		return fmt.Errorf("redis client not available")
	}

	lockKey := fmt.Sprintf("lock:%s", requestID)

	// Only release the lock if it belongs to this instance
	currentOwner, err := q.client.Get(ctx, lockKey)
	if err != nil {
		// Lock doesn't exist or expired, nothing to release
		return nil
	}

	if currentOwner == q.instanceID {
		err = q.client.Del(ctx, lockKey)
		if err != nil {
			return fmt.Errorf("failed to release lock for request %s: %w", requestID, err)
		}
		slog.Debug("Released distributed lock", "request_id", requestID, "instance_id", q.instanceID)
	}

	return nil
}

// SendHeartbeat updates the instance's health status in Redis
func (q *RedisQueue) SendHeartbeat(ctx context.Context) error {
	if q.client == nil {
		return fmt.Errorf("redis client not available")
	}

	q.lastHeartbeat = time.Now()
	heartbeatData := map[string]any{
		"instance_id": q.instanceID,
		"heartbeat":   q.lastHeartbeat.Unix(),
		"queue_size":  q.Size(),
		"processing":  q.processing,
	}

	data, err := json.Marshal(heartbeatData)
	if err != nil {
		return fmt.Errorf("failed to marshal heartbeat data: %w", err)
	}

	err = q.client.Set(ctx, q.healthKey, string(data), 60*time.Second) // 60 second TTL
	if err != nil {
		return fmt.Errorf("failed to send heartbeat: %w", err)
	}

	return nil
}

// GetInstanceID returns the instance ID (for testing purposes)
func (q *RedisQueue) GetInstanceID() string {
	return q.instanceID
}

// CleanupStaleLocks removes locks held by instances that are no longer active
func (q *RedisQueue) CleanupStaleLocks(ctx context.Context) error {
	if q.client == nil {
		return fmt.Errorf("redis client not available")
	}

	// Get all lock keys
	lockKeys, err := q.client.GetClient().Keys(ctx, "lock:*").Result()
	if err != nil {
		return fmt.Errorf("failed to get lock keys: %w", err)
	}

	cleaned := 0
	for _, lockKey := range lockKeys {
		// Extract request ID from lock key
		if len(lockKey) <= 5 { // "lock:" prefix
			continue
		}
		requestID := lockKey[5:] // Remove "lock:" prefix

		// Get the lock owner
		owner, err := q.client.Get(ctx, lockKey)
		if err != nil {
			continue // Lock might have expired
		}

		// Check if the owner instance is still active
		ownerHealthKey := fmt.Sprintf("instances:%s", owner)
		if _, err := q.client.Get(ctx, ownerHealthKey); err != nil {
			// Instance is not active, release the lock
			if err := q.client.Del(ctx, lockKey); err != nil {
				slog.Error("Failed to cleanup stale lock", "lock_key", lockKey, "error", err)
			} else {
				cleaned++
				slog.Info("Cleaned up stale lock", "request_id", requestID, "owner_instance", owner)
			}
		}
	}

	if cleaned > 0 {
		slog.Info("Cleaned up stale locks", "count", cleaned)
	}

	return nil
}
