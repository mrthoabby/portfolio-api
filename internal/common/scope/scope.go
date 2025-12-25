package scope

import "strings"

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

// IsProduction returns true if the scope is Production
func (s Scope) IsProduction() bool {
	return s == Production
}

// IsTest returns true if the scope is Test
func (s Scope) IsTest() bool {
	return s == Test
}

// IsLocal returns true if the scope is Local
func (s Scope) IsLocal() bool {
	return s == Local
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
