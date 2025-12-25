package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mrthoabby/portfolio-api/internal/common/scope"
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

	scope.Initialize()
	config, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "mongodb://localhost:27017", config.Database.URL)
	assert.Equal(t, "portfolio_db", config.Database.Name)
	assert.Equal(t, "8080", config.Server.Port) // default
	assert.True(t, scope.IsLocal())
}

func TestLoad_WithAllVariables(t *testing.T) {
	os.Clearenv()
	os.Setenv("PORT", "9090")
	os.Setenv("ENV", "production")
	os.Setenv("DATABASE_URL", "mongodb://prod:27017")
	os.Setenv("DATABASE_NAME", "prod_db")
	os.Setenv("ALLOWED_ORIGINS", "https://example.com,https://api.example.com")
	defer os.Clearenv()

	scope.Initialize()
	config, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "9090", config.Server.Port)
	assert.True(t, scope.IsProduction())
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

func TestParseScope(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected scope.Scope
	}{
		// Production variations
		{
			name:     "production lowercase",
			input:    "production",
			expected: scope.Production,
		},
		{
			name:     "production uppercase",
			input:    "PRODUCTION",
			expected: scope.Production,
		},
		{
			name:     "production mixed case",
			input:    "Production",
			expected: scope.Production,
		},
		{
			name:     "prod lowercase",
			input:    "prod",
			expected: scope.Production,
		},
		{
			name:     "prod uppercase",
			input:    "PROD",
			expected: scope.Production,
		},
		{
			name:     "prod mixed case",
			input:    "Prod",
			expected: scope.Production,
		},
		{
			name:     "production with prefix",
			input:    "my-production",
			expected: scope.Production,
		},
		{
			name:     "production with suffix",
			input:    "production-env",
			expected: scope.Production,
		},
		{
			name:     "prod with prefix",
			input:    "my-prod",
			expected: scope.Production,
		},
		{
			name:     "prod with suffix",
			input:    "prod-env",
			expected: scope.Production,
		},
		// Test variations
		{
			name:     "test lowercase",
			input:    "test",
			expected: scope.Test,
		},
		{
			name:     "test uppercase",
			input:    "TEST",
			expected: scope.Test,
		},
		{
			name:     "test mixed case",
			input:    "Test",
			expected: scope.Test,
		},
		{
			name:     "test with number",
			input:    "test-1",
			expected: scope.Test,
		},
		{
			name:     "test with alphanumeric",
			input:    "test-lfsj24",
			expected: scope.Test,
		},
		{
			name:     "test with prefix",
			input:    "my-test",
			expected: scope.Test,
		},
		{
			name:     "test with suffix",
			input:    "test-env",
			expected: scope.Test,
		},
		{
			name:     "test with underscore",
			input:    "test_env",
			expected: scope.Test,
		},
		// Local/default variations
		{
			name:     "development",
			input:    "development",
			expected: scope.Local,
		},
		{
			name:     "local",
			input:    "local",
			expected: scope.Local,
		},
		{
			name:     "empty string",
			input:    "",
			expected: scope.Local,
		},
		{
			name:     "staging",
			input:    "staging",
			expected: scope.Local,
		},
		{
			name:     "dev",
			input:    "dev",
			expected: scope.Local,
		},
		{
			name:     "unknown value",
			input:    "unknown",
			expected: scope.Local,
		},
		{
			name:     "Production always is the first priority when the value is present and Test is not",
			input:    "production-test",
			expected: scope.Production,
		},
		{
			name:     "Production is always the first priority: When is present does not matter the order",
			input:    "test-production",
			expected: scope.Production,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scope.ParseScope(tt.input)
			assert.Equal(t, tt.expected, result, "ParseScope(%q) = %v, want %v", tt.input, result, tt.expected)
		})
	}
}

func TestScope_IsProduction(t *testing.T) {
	os.Clearenv()
	os.Setenv("ENV", "production")
	scope.Initialize()
	assert.True(t, scope.IsProduction())
	assert.False(t, scope.IsTest())
	assert.False(t, scope.IsLocal())
	defer os.Clearenv()
}

func TestScope_IsTest(t *testing.T) {
	os.Clearenv()
	os.Setenv("ENV", "test")
	scope.Initialize()
	assert.False(t, scope.IsProduction())
	assert.True(t, scope.IsTest())
	assert.False(t, scope.IsLocal())
	defer os.Clearenv()
}

func TestScope_IsLocal(t *testing.T) {
	os.Clearenv()
	os.Setenv("ENV", "development")
	scope.Initialize()
	assert.False(t, scope.IsProduction())
	assert.False(t, scope.IsTest())
	assert.True(t, scope.IsLocal())
	defer os.Clearenv()
}

func TestScope_String(t *testing.T) {
	os.Clearenv()
	os.Setenv("ENV", "production")
	scope.Initialize()
	assert.Equal(t, "production", scope.String())
	defer os.Clearenv()

	os.Clearenv()
	os.Setenv("ENV", "test")
	scope.Initialize()
	assert.Equal(t, "test", scope.String())
	defer os.Clearenv()

	os.Clearenv()
	os.Setenv("ENV", "development")
	scope.Initialize()
	assert.Equal(t, "local", scope.String())
	defer os.Clearenv()
}

func TestLoad_WithTestEnvironment(t *testing.T) {
	os.Clearenv()
	os.Setenv("PORT", "8080")
	os.Setenv("ENV", "test-1")
	os.Setenv("DATABASE_URL", "mongodb://test:27017")
	os.Setenv("DATABASE_NAME", "test_db")
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	defer os.Clearenv()

	scope.Initialize()
	config, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.True(t, scope.IsTest())
	assert.False(t, scope.IsProduction())
	assert.False(t, scope.IsLocal())
}

func TestLoad_WithProdVariation(t *testing.T) {
	os.Clearenv()
	os.Setenv("PORT", "8080")
	os.Setenv("ENV", "prod") // Using "prod" instead of "production"
	os.Setenv("DATABASE_URL", "mongodb://prod:27017")
	os.Setenv("DATABASE_NAME", "prod_db")
	os.Setenv("ALLOWED_ORIGINS", "https://example.com")
	defer os.Clearenv()

	scope.Initialize()
	config, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.True(t, scope.IsProduction())
}
