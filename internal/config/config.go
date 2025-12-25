package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/mrthoabby/portfolio-api/internal/common/logger"
	"github.com/mrthoabby/portfolio-api/internal/common/scope"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL  string
	Name string
}

type CORSConfig struct {
	AllowedOrigins []string
}

// Load loads configuration from environment variables.
// Required variables: DATABASE_URL, DATABASE_NAME, ALLOWED_ORIGINS, PORT
// This function uses the standard log package for backward compatibility.
func Load() (*Config, error) {
	return LoadWithLogger(nil)
}

// LoadWithLogger loads configuration from environment variables using the provided logger.
// If logger is nil, it falls back to the standard log package.
// Required variables: DATABASE_URL, DATABASE_NAME, ALLOWED_ORIGINS, PORT
func LoadWithLogger(appLogger logger.Logger) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		if appLogger != nil {
			appLogger.Debug(".env file not found or error loading, using environment variables only",
				logger.Error(err),
			)
		} else {
			log.Printf("Config: .env file not found or error loading: %v (using environment variables only)", err)
		}
	} else {
		if appLogger != nil {
			appLogger.Debug(".env file loaded successfully")
		} else {
			log.Println("Config: .env file loaded successfully")
		}
	}

	// Validate required environment variables
	var missingVars []string

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		missingVars = append(missingVars, "DATABASE_URL")
	}

	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		missingVars = append(missingVars, "DATABASE_NAME")
	} else {
		if appLogger != nil {
			appLogger.Debug("DATABASE_NAME loaded",
				logger.String("database", databaseName),
			)
		} else {
			log.Printf("Config: DATABASE_NAME loaded: %s", databaseName)
		}
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		missingVars = append(missingVars, "ALLOWED_ORIGINS")
	} else {
		if appLogger != nil {
			appLogger.Debug("ALLOWED_ORIGINS loaded",
				logger.String("origins", allowedOrigins),
			)
		} else {
			log.Printf("Config: ALLOWED_ORIGINS loaded: %s", allowedOrigins)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		missingVars = append(missingVars, "PORT")
	} else {
		if appLogger != nil {
			appLogger.Debug("PORT loaded",
				logger.String("port", port),
			)
		} else {
			log.Printf("Config: PORT loaded: %s", port)
		}
	}

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missingVars, ", "))
	}

	config := &Config{
		Server: ServerConfig{
			Port: port,
		},
		Database: DatabaseConfig{
			URL:  databaseURL,
			Name: databaseName,
		},
		CORS: CORSConfig{
			AllowedOrigins: parseOrigins(allowedOrigins),
		},
	}

	if appLogger != nil {
		appLogger.Info("Configuration loaded successfully",
			logger.String("scope", scope.String()),
			logger.String("port", port),
		)
	} else {
		log.Printf("Config: Configuration loaded successfully (scope: %s, port: %s)", scope.String(), port)
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseOrigins(origins string) []string {
	if origins == "" {
		return []string{}
	}
	return strings.Split(origins, ",")
}
