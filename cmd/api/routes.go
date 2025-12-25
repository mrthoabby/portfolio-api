package main

import (
	"github.com/go-chi/chi/v5"

	"github.com/mrthoabby/portfolio-api/internal/common/logger"
	"github.com/mrthoabby/portfolio-api/internal/middleware"
)

// SetupRoutes configures all routes and middleware
func SetupRoutes(deps *Dependencies, appLogger logger.Logger, allowedOrigins []string) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware (order matters!)
	r.Use(middleware.RecoverPanic)             // Recover from panics first
	r.Use(middleware.RequestID)                // Add request ID for tracing
	r.Use(middleware.ClientIP)                 // Extract client IP to context
	r.Use(middleware.WithLogger(appLogger))    // Log requests with structured logger
	r.Use(middleware.SecurityHeaders)          // Add security headers
	r.Use(middleware.MaxBodySize(maxBodySize)) // Limit request body size
	r.Use(deps.GlobalRateLimiter.Limit)        // Global rate limiting
	r.Use(middleware.NewCORS(allowedOrigins))  // CORS

	// Health check (no rate limiting needed)
	r.Get("/health", deps.HealthHandler.Check)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/profiles/{id}", func(r chi.Router) {
			r.Get("/", deps.ProfileHandler.GetByID)
			r.Get("/skills", deps.SkillsHandler.GetByProfileID)
			r.Get("/projects", deps.ProjectsHandler.GetByProfileID)
			r.Get("/certificates", deps.CertificatesHandler.GetByProfileID)

			// Contact endpoint with stricter rate limiting
			r.With(deps.ContactRateLimiter.Limit).Post("/contacts", deps.ContactsHandler.Create)

			// Questions endpoint with rate limiting
			r.With(deps.QuestionRateLimiter.Limit).Post("/questions", deps.QuestionsHandler.Create)
		})
	})

	return r
}
