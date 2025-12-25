package main

import (
	"os"
	"strings"

	"github.com/mrthoabby/portfolio-api/internal/common/logger"
	"github.com/mrthoabby/portfolio-api/internal/common/scope"
	"github.com/mrthoabby/portfolio-api/internal/config"
	"github.com/mrthoabby/portfolio-api/internal/repository/mongo"
	"github.com/mrthoabby/portfolio-api/internal/version"
)

func main() {
	// Initialize scope first
	scope.Initialize()

	// Initialize logger (uses the initialized scope)
	appLogger := logger.MustNewLogger()
	defer func() {
		if syncLogger, ok := appLogger.(interface{ Sync() error }); ok {
			_ = syncLogger.Sync()
		}
	}()

	// Log version information
	v := version.Get()
	appLogger.Info("Starting Portfolio API",
		logger.String("version", v.Version),
		logger.String("built", v.BuildDate),
	)

	// Load configuration
	appLogger.Info("Loading configuration...")
	cfg, err := config.LoadWithLogger(appLogger)
	if err != nil {
		appLogger.Error("Failed to load config", logger.Error(err))
		os.Exit(1)
	}

	// Validate database URL is not localhost in production
	if scope.IsProduction() {
		if strings.Contains(cfg.Database.URL, "localhost") || strings.Contains(cfg.Database.URL, "127.0.0.1") {
			appLogger.Error("FATAL: DATABASE_URL contains localhost/127.0.0.1 in production environment",
				logger.String("message", "This is not allowed. Please use a remote MongoDB instance (e.g., MongoDB Atlas)."),
			)
			os.Exit(1)
		}
	}

	// Connect to data source (MongoDB)
	appLogger.Info("Connecting to MongoDB database",
		logger.String("database", cfg.Database.Name),
	)
	dataSource, err := mongo.NewDataSource(cfg.Database.URL, cfg.Database.Name, appLogger)
	if err != nil {
		appLogger.Error("Failed to connect to data source", logger.Error(err))
		os.Exit(1)
	}
	defer dataSource.Close()

	appLogger.Info("Data source connection established successfully")

	// Initialize dependencies
	deps := InitializeDependencies(dataSource, appLogger)

	// Setup routes
	router := SetupRoutes(deps, appLogger, cfg.CORS.AllowedOrigins)

	// Create and start server
	server := NewServer(cfg, router, appLogger)
	if err := server.Start(); err != nil {
		appLogger.Error("Failed to start server", logger.Error(err))
		os.Exit(1)
	}

	// Wait for shutdown signal
	server.WaitForShutdown()
}
