package queue

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryQueue_EnqueueDequeue(t *testing.T) {
	q := NewInMemoryQueue()
	ctx := context.Background()

	// Test enqueue
	req := &ExecutionRequest{
		RuleID:    uuid.New(),
		TriggerID: uuid.New(),
		EventData: map[string]interface{}{"test": "data"},
	}

	err := q.Enqueue(ctx, req)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, req.ID)
	assert.False(t, req.QueuedAt.IsZero())

	// Test dequeue
	dequeued, err := q.Dequeue(ctx)
	require.NoError(t, err)
	assert.Equal(t, req.ID, dequeued.ID)
	assert.Equal(t, req.RuleID, dequeued.RuleID)
	assert.Equal(t, req.TriggerID, dequeued.TriggerID)
	assert.Equal(t, req.EventData, dequeued.EventData)
}

func TestInMemoryQueue_Size(t *testing.T) {
	q := NewInMemoryQueue()
	ctx := context.Background()

	// Initially empty
	assert.Equal(t, 0, q.Size())

	// Add items
	req1 := &ExecutionRequest{RuleID: uuid.New()}
	req2 := &ExecutionRequest{RuleID: uuid.New()}

	q.Enqueue(ctx, req1)
	assert.Equal(t, 1, q.Size())

	q.Enqueue(ctx, req2)
	assert.Equal(t, 2, q.Size())

	// Remove items
	q.Dequeue(ctx)
	assert.Equal(t, 1, q.Size())

	q.Dequeue(ctx)
	assert.Equal(t, 0, q.Size())
}

func TestInMemoryQueue_DequeueEmpty(t *testing.T) {
	q := NewInMemoryQueue()
	ctx := context.Background()

	_, err := q.Dequeue(ctx)
	assert.Equal(t, ErrQueueEmpty, err)
}

func TestInMemoryQueue_EnqueueClosed(t *testing.T) {
	q := NewInMemoryQueue()
	ctx := context.Background()

	q.Close()

	req := &ExecutionRequest{RuleID: uuid.New()}
	err := q.Enqueue(ctx, req)
	assert.Equal(t, ErrQueueClosed, err)
}

func TestInMemoryQueue_DequeueClosed(t *testing.T) {
	q := NewInMemoryQueue()
	ctx := context.Background()

	q.Close()

	_, err := q.Dequeue(ctx)
	assert.Equal(t, ErrQueueClosed, err)
}

func TestInMemoryQueue_Close(t *testing.T) {
	q := NewInMemoryQueue()
	ctx := context.Background()

	// Add an item
	req := &ExecutionRequest{RuleID: uuid.New()}
	q.Enqueue(ctx, req)

	// Close the queue
	err := q.Close()
	require.NoError(t, err)

	// Should not be able to enqueue after close
	err = q.Enqueue(ctx, &ExecutionRequest{RuleID: uuid.New()})
	assert.Equal(t, ErrQueueClosed, err)

	// Should still be able to dequeue existing items
	dequeued, err := q.Dequeue(ctx)
	require.NoError(t, err)
	assert.Equal(t, req.RuleID, dequeued.RuleID)

	// Should get closed error after queue is empty
	_, err = q.Dequeue(ctx)
	assert.Equal(t, ErrQueueClosed, err)
}

func TestRedisQueue_BasicOperations(t *testing.T) {
	// Skip if Redis is not available
	if testing.Short() {
		t.Skip("Skipping Redis queue test in short mode")
	}

	// This test would require a Redis instance
	// For now, just test that the constructor works with nil client
	q := NewRedisQueue(nil, "test:queue")
	assert.NotNil(t, q)

	// Should return error for operations without Redis client
	ctx := context.Background()
	req := &ExecutionRequest{RuleID: uuid.New()}

	err := q.Enqueue(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Redis client not available")

	_, err = q.Dequeue(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Redis client not available")

	size := q.Size()
	assert.Equal(t, 0, size)

	err = q.Close()
	assert.NoError(t, err)
}

func TestRedisQueue_DistributedLocking(t *testing.T) {
	// Skip if Redis is not available
	if testing.Short() {
		t.Skip("Skipping Redis distributed locking test in short mode")
	}

	// This test would require a Redis instance
	// For now, just test that the constructor initializes instance ID
	q := NewRedisQueue(nil, "test:queue")
	assert.NotNil(t, q)
	assert.NotEmpty(t, q.GetInstanceID())
	assert.Contains(t, q.GetInstanceID(), "-") // Should contain hostname-pid-random format
}

func TestGenerateInstanceID(t *testing.T) {
	id1 := generateInstanceID()
	id2 := generateInstanceID()

	// IDs should be unique
	assert.NotEqual(t, id1, id2)

	// Should contain hostname, PID, and random component
	assert.Contains(t, id1, "-")
	assert.Contains(t, id2, "-")

	// Should be reasonably long (hostname + pid + random)
	assert.Greater(t, len(id1), 10)
	assert.Greater(t, len(id2), 10)
}
