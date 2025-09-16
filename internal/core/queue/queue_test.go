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
