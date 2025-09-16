package tracing

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// TracerProvider holds the OpenTelemetry tracer provider
var TracerProvider *sdktrace.TracerProvider

// Tracer is the global tracer instance
var Tracer trace.Tracer

// InitTracing initializes OpenTelemetry tracing
func InitTracing(serviceName, serviceVersion string) error {
	// Create resource with service information
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Create stdout exporter for development
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return fmt.Errorf("failed to create stdout trace exporter: %w", err)
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)
	TracerProvider = tp

	// Create and set global tracer
	Tracer = tp.Tracer(serviceName)

	log.Printf("OpenTelemetry tracing initialized for service: %s v%s", serviceName, serviceVersion)
	return nil
}

// ShutdownTracing shuts down the tracer provider
func ShutdownTracing(ctx context.Context) error {
	if TracerProvider != nil {
		return TracerProvider.Shutdown(ctx)
	}
	return nil
}

// StartSpan starts a new span with the given name
func StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	if Tracer == nil {
		// Tracing not initialized, return no-op span
		return ctx, trace.SpanFromContext(ctx)
	}
	return Tracer.Start(ctx, spanName)
}

// StartSpanFromContext starts a new span from the given context
func StartSpanFromContext(ctx context.Context, spanName string) (context.Context, trace.Span) {
	if Tracer == nil {
		// Tracing not initialized, return no-op span
		return ctx, trace.SpanFromContext(ctx)
	}
	return Tracer.Start(ctx, spanName)
}
