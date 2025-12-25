package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/mrthoabby/portfolio-api/internal/common"
	"github.com/mrthoabby/portfolio-api/internal/common/logger"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (responseWriter *responseWriter) WriteHeader(code int) {
	responseWriter.statusCode = code
	responseWriter.ResponseWriter.WriteHeader(code)
}

// Logger creates a middleware that logs HTTP requests using the standard log package.
// For structured logging, use WithLogger instead.
func Logger(next http.Handler) http.Handler {
	return WithLogger(nil)(next)
}

// WithLogger returns a middleware function that logs HTTP requests using the provided logger.
// If logger is nil, it falls back to the standard log package.
func WithLogger(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			requestID := GetRequestID(r.Context())
			clientIP := common.ClientIPFromContext(r.Context())

			if l != nil {
				// Structured logging with fields
				l.Info("HTTP request",
					logger.String("request_id", requestID),
					logger.String("method", r.Method),
					logger.String("path", r.URL.Path),
					logger.Int("status", wrapped.statusCode),
					logger.Duration("duration", duration),
					logger.String("client_ip", clientIP),
				)
			} else {
				// Fallback to standard log format
				log.Printf(
					"[%s] %s %s %d %v %s",
					requestID,
					r.Method,
					r.URL.Path,
					wrapped.statusCode,
					duration,
					clientIP,
				)
			}
		})
	}
}
