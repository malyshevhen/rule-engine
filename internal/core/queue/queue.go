package queue

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/metrics"
)

// Common queue errors
var (
	ErrQueueClosed = errors.New("queue is closed")
	ErrQueueEmpty  = errors.New("queue is empty")
)

// ExecutionRequest represents a rule execution request
type ExecutionRequest struct {
	ID        uuid.UUID              `json:"id"`
	RuleID    uuid.UUID              `json:"rule_id"`
	TriggerID uuid.UUID              `json:"trigger_id"`
	EventData map[string]interface{} `json:"event_data,omitempty"`
	QueuedAt  time.Time              `json:"queued_at"`
}

// Queue interface for rule execution queuing
type Queue interface {
	Enqueue(ctx context.Context, req *ExecutionRequest) error
	Dequeue(ctx context.Context) (*ExecutionRequest, error)
	Size() int
	Close() error
}

// InMemoryQueue implements an in-memory queue for rule executions
type InMemoryQueue struct {
	mu     sync.RWMutex
	queue  []*ExecutionRequest
	closed bool
}

// NewInMemoryQueue creates a new in-memory queue
func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{
		queue: make([]*ExecutionRequest, 0),
	}
}

// Enqueue adds a rule execution request to the queue
func (q *InMemoryQueue) Enqueue(ctx context.Context, req *ExecutionRequest) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return ErrQueueClosed
	}

	req.ID = uuid.New()
	req.QueuedAt = time.Now()

	q.queue = append(q.queue, req)

	// Update metrics
	metrics.QueueSize.Set(float64(len(q.queue)))
	metrics.QueueEnqueueTotal.Inc()

	return nil
}

// Dequeue removes and returns the next rule execution request from the queue
func (q *InMemoryQueue) Dequeue(ctx context.Context) (*ExecutionRequest, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queue) == 0 {
		if q.closed {
			return nil, ErrQueueClosed
		}
		return nil, ErrQueueEmpty
	}

	req := q.queue[0]
	q.queue = q.queue[1:]

	// Update metrics
	metrics.QueueSize.Set(float64(len(q.queue)))
	metrics.QueueDequeueTotal.Inc()

	return req, nil
}

// Size returns the current number of items in the queue
func (q *InMemoryQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.queue)
}

// Close closes the queue
func (q *InMemoryQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
	return nil
}
