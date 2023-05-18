package utils

import (
	"context"

	telemetry "github.com/bancodobrasil/gin-telemetry"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
)

// GenerateSpanTracer generates a span tracer for telemetry purposes, but only if middleware is enable.
func GenerateSpanTracer(ctx context.Context, name string) func() {
	MiddlewareDisabled := viper.GetViper().GetBool("TELEMETRY_DISABLED")
	if !MiddlewareDisabled {
		tracer := telemetry.GetTracer(ctx)
		_, span := tracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindInternal))
		return func() {
			span.End()
		}
	}

	return func() {}
}
