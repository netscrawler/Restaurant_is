package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// responseBodyWriter перехватывает и сохраняет body ответа
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggingMiddleware возвращает gin.HandlerFunc, который логирует HTTP-запросы и ответы с trace_id и всей информацией о запросе и ответе.
func LoggingMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceID := trace.SpanContextFromContext(c.Request.Context()).TraceID().String()

		// Чтение тела запроса
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			// Восстановить тело для дальнейшего использования
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		log.Info("HTTP REQ",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("trace_id", traceID),
			slog.String("query", c.Request.URL.RawQuery),
			slog.Any("headers", c.Request.Header),
			slog.String("body", string(bodyBytes)),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.String("host", c.Request.Host),
			slog.String("proto", c.Request.Proto),
			slog.Int64("content_length", c.Request.ContentLength),
		)

		// Подмена ResponseWriter для перехвата body ответа
		rw := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = rw

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Info("HTTP RESP",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", status),
			slog.String("latency", latency.String()),
			slog.String("trace_id", traceID),
			slog.Any("headers", c.Writer.Header()),
			slog.Int("size", c.Writer.Size()),
			slog.String("content_type", c.Writer.Header().Get("Content-Type")),
			slog.Any("cookies", c.Writer.Header()["Set-Cookie"]),
			slog.String("body", rw.body.String()),
		)
	}
}
