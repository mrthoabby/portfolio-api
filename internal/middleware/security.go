package middleware

import (
	"net/http"

	"github.com/mrthoabby/portfolio-api/internal/common"
)

// SecurityHeaders adds security-related HTTP headers to responses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, r *http.Request) {
		// Prevent MIME type sniffing
		responseWriter.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		responseWriter.Header().Set("X-Frame-Options", "DENY")

		// XSS protection (legacy but still useful)
		responseWriter.Header().Set("X-XSS-Protection", "1; mode=block")

		// Referrer policy
		responseWriter.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy (restrictive for API)
		responseWriter.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")

		// Permissions Policy (disable unnecessary features)
		responseWriter.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Cache control for API responses (public data, cacheable for 5 minutes)
		responseWriter.Header().Set("Cache-Control", "public, max-age=3000")

		next.ServeHTTP(responseWriter, r)
	})
}

// MaxBodySize limits the size of request bodies
func MaxBodySize(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only limit POST, PUT, PATCH requests
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RecoverPanic recovers from panics and returns a 500 error
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic (in production, use proper logging)
				common.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
