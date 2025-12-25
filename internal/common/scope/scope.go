package scope

import (
	"os"
	"strings"
)

// Scope represents the environment scope of the application
type Scope int

const (
	// Local is the default scope for local development
	Local Scope = iota
	// Production represents production environments
	Production
	// Test represents test/testing environments
	Test
)

// currentScope holds the current application scope
var currentScope Scope

// Initialize initializes the scope from the ENV environment variable.
// This should be called once at application startup.
func Initialize() {
	currentScope = ParseScope(os.Getenv("ENV"))
}

// String returns the string representation of the scope
func (s Scope) String() string {
	switch s {
	case Production:
		return "production"
	case Test:
		return "test"
	default:
		return "local"
	}
}

// ParseScope parses a string value and returns the corresponding Scope.
// Detection logic:
// - Production: contains "prod" or "production" (case-insensitive)
// - Test: contains "test" (case-insensitive)
// - Local: everything else (default)
func ParseScope(envValue string) Scope {
	if envValue == "" {
		return Local
	}

	lower := strings.ToLower(envValue)

	// Check for production: contains "prod" or "production"
	if strings.Contains(lower, "prod") || strings.Contains(lower, "production") {
		return Production
	}

	// Check for test: contains "test"
	if strings.Contains(lower, "test") {
		return Test
	}

	// Default to local
	return Local
}

// IsProduction returns true if the current scope is Production
func IsProduction() bool {
	return currentScope == Production
}

// IsTest returns true if the current scope is Test
func IsTest() bool {
	return currentScope == Test
}

// IsLocal returns true if the current scope is Local
func IsLocal() bool {
	return currentScope == Local
}

// String returns the string representation of the current scope
func String() string {
	return currentScope.String()
}
