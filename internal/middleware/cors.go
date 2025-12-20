package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

func NewCORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: false, // Don't allow credentials for public API
		MaxAge:           3600,  // Cache preflight for 1 hour
	})
}
