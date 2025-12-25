package main

import (
	"github.com/mrthoabby/portfolio-api/internal/application/certificates"
	"github.com/mrthoabby/portfolio-api/internal/application/contacts"
	"github.com/mrthoabby/portfolio-api/internal/application/profile"
	"github.com/mrthoabby/portfolio-api/internal/application/projects"
	"github.com/mrthoabby/portfolio-api/internal/application/questions"
	"github.com/mrthoabby/portfolio-api/internal/application/skills"
	"github.com/mrthoabby/portfolio-api/internal/common/contracts"
	"github.com/mrthoabby/portfolio-api/internal/common/logger"
	"github.com/mrthoabby/portfolio-api/internal/health"
	"github.com/mrthoabby/portfolio-api/internal/middleware"
)

// Dependencies holds all application dependencies
type Dependencies struct {
	// Rate limiters
	GlobalRateLimiter   *middleware.RateLimiter
	ContactRateLimiter  *middleware.RateLimiter
	QuestionRateLimiter *middleware.RateLimiter

	// Handlers
	ProfileHandler      *profile.Handler
	SkillsHandler       *skills.Handler
	ProjectsHandler     *projects.Handler
	CertificatesHandler *certificates.Handler
	ContactsHandler     *contacts.Handler
	QuestionsHandler    *questions.Handler
	HealthHandler       *health.Handler
}

// InitializeDependencies initializes all application dependencies
func InitializeDependencies(dataSource contracts.DataSource, appLogger logger.Logger) *Dependencies {
	// Initialize rate limiters
	globalRateLimiter := middleware.NewRateLimiter(rateLimitRequests, rateLimitWindow)
	contactRateLimiter := middleware.NewRateLimiter(contactRateLimit, contactRateWindow)
	questionRateLimiter := middleware.NewRateLimiter(questionRateLimit, questionRateWindow)

	// Initialize profile domain
	profileRepo := profile.NewRepository(dataSource)
	profileService := profile.NewService(profileRepo)
	profileHandler := profile.NewHandler(profileService)

	// Initialize skills domain
	skillsRepo := skills.NewRepository(dataSource)
	skillsService := skills.NewService(skillsRepo, profileService)
	skillsHandler := skills.NewHandler(skillsService)

	// Initialize projects domain
	projectsRepo := projects.NewRepository(dataSource)
	projectsService := projects.NewService(projectsRepo, profileService)
	projectsHandler := projects.NewHandler(projectsService)

	// Initialize certificates domain
	certificatesRepo := certificates.NewRepository(dataSource)
	certificatesService := certificates.NewService(certificatesRepo, profileService)
	certificatesHandler := certificates.NewHandler(certificatesService)

	// Initialize contacts domain
	contactsRepo := contacts.NewRepository(dataSource)
	contactsService := contacts.NewService(contactsRepo, profileService)
	contactsHandler := contacts.NewHandler(contactsService)

	// Initialize questions domain
	questionsRepo := questions.NewRepository(dataSource)
	questionsService := questions.NewService(questionsRepo, profileService)
	questionsHandler := questions.NewHandler(questionsService)

	// Initialize health handler
	healthHandler := health.NewHandler(dataSource)

	return &Dependencies{
		GlobalRateLimiter:   globalRateLimiter,
		ContactRateLimiter:  contactRateLimiter,
		QuestionRateLimiter: questionRateLimiter,
		ProfileHandler:      profileHandler,
		SkillsHandler:       skillsHandler,
		ProjectsHandler:     projectsHandler,
		CertificatesHandler: certificatesHandler,
		ContactsHandler:     contactsHandler,
		QuestionsHandler:    questionsHandler,
		HealthHandler:       healthHandler,
	}
}
