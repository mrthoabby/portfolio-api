package logger

import "fmt"

// NewLogger creates a new logger instance based on the current scope.
// The scope must be initialized with scope.Initialize() before calling this function.
// - Production scope -> JSON format, Info level
// - anything else -> Development format (text), Debug level
//
// Returns a Logger interface that can be used throughout the application.
func NewLogger() (Logger, error) {
	return NewZapLogger()
}

// MustNewLogger creates a new logger and panics if it fails.
// Use this in main() or initialization code where logger failure is fatal.
func MustNewLogger() Logger {
	logger, err := NewLogger()
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %v", err))
	}
	return logger
}
