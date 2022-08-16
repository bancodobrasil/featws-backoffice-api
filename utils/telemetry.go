package utils

import (
	"context"

	telemetry "github.com/bancodobrasil/gin-telemetry"
	"go.opentelemetry.io/otel/trace"
)

// GenerateSpanTracer ...
func GenerateSpanTracer(ctx context.Context, name string) func() {
	tracer := telemetry.GetTracer(ctx)
	_, span := tracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindInternal))
	return func() {
		span.End()
	}
}
