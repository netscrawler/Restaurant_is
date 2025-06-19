package utils

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

// LoggerWithTrace возвращает логгер с trace_id из контекста.
func LoggerWithTrace(ctx context.Context, baseLogger *slog.Logger) *slog.Logger {
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()

	return baseLogger.With(slog.String("trace_id", traceID))
}
