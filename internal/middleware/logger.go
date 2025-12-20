package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (responseWriter *responseWriter) WriteHeader(code int) {
	responseWriter.statusCode = code
	responseWriter.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		requestID := GetRequestID(r.Context())

		// Structured log format: [REQUEST_ID] METHOD PATH STATUS DURATION CLIENT_IP
		log.Printf(
			"[%s] %s %s %d %v %s",
			requestID,
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
			getIP(r),
		)
	})
}
