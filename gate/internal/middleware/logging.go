package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// LoggingMiddleware возвращает gin.HandlerFunc, который логирует HTTP-запросы и ответы с trace_id.
func LoggingMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceID := trace.SpanContextFromContext(c.Request.Context()).TraceID().String()

		log.Info("HTTP REQ",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("trace_id", traceID),
		)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Info("HTTP RESP",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", status),
			slog.String("latency", latency.String()),
			slog.String("trace_id", traceID),
		)
	}
}
