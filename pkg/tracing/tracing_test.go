package tracing

import (
	"context"
	"testing"
)

func TestInitTracing(t *testing.T) {
	// Test initialization
	err := InitTracing("test-service", "1.0.0")
	if err != nil {
		t.Fatalf("InitTracing failed: %v", err)
	}

	if TracerProvider == nil {
		t.Error("TracerProvider should not be nil after initialization")
	}

	if Tracer == nil {
		t.Error("Tracer should not be nil after initialization")
	}

	// Test shutdown
	ctx := context.Background()
	err = ShutdownTracing(ctx)
	if err != nil {
		t.Fatalf("ShutdownTracing failed: %v", err)
	}
}

func TestStartSpan(t *testing.T) {
	// Initialize tracing
	err := InitTracing("test-service", "1.0.0")
	if err != nil {
		t.Fatalf("InitTracing failed: %v", err)
	}
	defer func() {
		_ = ShutdownTracing(context.Background())
	}()

	ctx := context.Background()
	newCtx, span := StartSpan(ctx, "test-span")

	if newCtx == nil {
		t.Error("New context should not be nil")
	}

	if span == nil {
		t.Error("Span should not be nil")
	}

	span.End()
}

func TestStartSpanFromContext(t *testing.T) {
	// Initialize tracing
	err := InitTracing("test-service", "1.0.0")
	if err != nil {
		t.Fatalf("InitTracing failed: %v", err)
	}
	defer func() {
		_ = ShutdownTracing(context.Background())
	}()

	ctx := context.Background()
	newCtx, span := StartSpanFromContext(ctx, "test-span")

	if newCtx == nil {
		t.Error("New context should not be nil")
	}

	if span == nil {
		t.Error("Span should not be nil")
	}

	span.End()
}

func TestStartSpanWithoutInit(t *testing.T) {
	// Reset globals for test
	TracerProvider = nil
	Tracer = nil

	ctx := context.Background()
	newCtx, span := StartSpan(ctx, "test-span")

	// Should return no-op span
	if newCtx != ctx {
		t.Error("Context should be unchanged when tracing not initialized")
	}

	if span.SpanContext().IsValid() {
		t.Error("Span should be no-op when tracing not initialized")
	}
}

func TestStartSpanFromContextWithoutInit(t *testing.T) {
	// Reset globals for test
	TracerProvider = nil
	Tracer = nil

	ctx := context.Background()
	newCtx, span := StartSpanFromContext(ctx, "test-span")

	// Should return no-op span
	if newCtx != ctx {
		t.Error("Context should be unchanged when tracing not initialized")
	}

	if span.SpanContext().IsValid() {
		t.Error("Span should be no-op when tracing not initialized")
	}
}
