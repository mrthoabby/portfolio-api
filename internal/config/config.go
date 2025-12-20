package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	URL  string
	Name string
}

type CORSConfig struct {
	AllowedOrigins []string
}

// Load loads configuration from environment variables.
// Required variables: DATABASE_URL, DATABASE_NAME, ALLOWED_ORIGINS
// Optional variables: PORT (default: 8080), ENV (default: development)
func Load() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	// Validate required environment variables
	var missingVars []string

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		missingVars = append(missingVars, "DATABASE_URL")
	}

	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		missingVars = append(missingVars, "DATABASE_NAME")
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		missingVars = append(missingVars, "ALLOWED_ORIGINS")
	}

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missingVars, ", "))
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("PORT", "8080"),
			Env:  getEnvOrDefault("ENV", "development"),
		},
		Database: DatabaseConfig{
			URL:  databaseURL,
			Name: databaseName,
		},
		CORS: CORSConfig{
			AllowedOrigins: parseOrigins(allowedOrigins),
		},
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
