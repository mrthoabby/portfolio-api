package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_MissingRequiredVariables(t *testing.T) {
	// Clear all environment variables
	os.Clearenv()

	config, err := Load()
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "DATABASE_URL")
	assert.Contains(t, err.Error(), "DATABASE_NAME")
	assert.Contains(t, err.Error(), "ALLOWED_ORIGINS")
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	os.Clearenv()
	os.Setenv("DATABASE_NAME", "test_db")
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	defer os.Clearenv()

	config, err := Load()
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "DATABASE_URL")
}

func TestLoad_MissingDatabaseName(t *testing.T) {
	os.Clearenv()
	os.Setenv("DATABASE_URL", "mongodb://localhost:27017")
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	defer os.Clearenv()

	config, err := Load()
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "DATABASE_NAME")
}

func TestLoad_MissingAllowedOrigins(t *testing.T) {
	os.Clearenv()
	os.Setenv("DATABASE_URL", "mongodb://localhost:27017")
	os.Setenv("DATABASE_NAME", "test_db")
	defer os.Clearenv()

	config, err := Load()
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "ALLOWED_ORIGINS")
}

func TestLoad_WithAllRequiredVariables(t *testing.T) {
	os.Clearenv()
	os.Setenv("DATABASE_URL", "mongodb://localhost:27017")
	os.Setenv("DATABASE_NAME", "portfolio_db")
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	defer os.Clearenv()

	config, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "mongodb://localhost:27017", config.Database.URL)
	assert.Equal(t, "portfolio_db", config.Database.Name)
	assert.Equal(t, "8080", config.Server.Port) // default
	assert.Equal(t, "development", config.Server.Env) // default
}

func TestLoad_WithAllVariables(t *testing.T) {
	os.Clearenv()
	os.Setenv("PORT", "9090")
	os.Setenv("ENV", "production")
	os.Setenv("DATABASE_URL", "mongodb://prod:27017")
	os.Setenv("DATABASE_NAME", "prod_db")
	os.Setenv("ALLOWED_ORIGINS", "https://example.com,https://api.example.com")
	defer os.Clearenv()

	config, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "9090", config.Server.Port)
	assert.Equal(t, "production", config.Server.Env)
	assert.Equal(t, "mongodb://prod:27017", config.Database.URL)
	assert.Equal(t, "prod_db", config.Database.Name)
	assert.Len(t, config.CORS.AllowedOrigins, 2)
	assert.Contains(t, config.CORS.AllowedOrigins, "https://example.com")
	assert.Contains(t, config.CORS.AllowedOrigins, "https://api.example.com")
}

func TestParseOrigins(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "multiple origins",
			input:    "http://localhost:3000,http://localhost:3001",
			expected: []string{"http://localhost:3000", "http://localhost:3001"},
		},
		{
			name:     "single origin",
			input:    "http://localhost:3000",
			expected: []string{"http://localhost:3000"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOrigins(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	os.Setenv("TEST_KEY", "test-value")
	defer os.Unsetenv("TEST_KEY")

	assert.Equal(t, "test-value", getEnvOrDefault("TEST_KEY", "default"))
	assert.Equal(t, "default", getEnvOrDefault("NON_EXISTENT", "default"))
}
