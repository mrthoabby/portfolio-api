package logger

import "context"

// Logger defines the interface for structured logging.
// This abstraction allows changing the logging implementation
// without modifying the business logic code.
type Logger interface {
	// Debug logs a debug message with optional fields
	Debug(msg string, fields ...Field)

	// Info logs an informational message with optional fields
	Info(msg string, fields ...Field)

	// Warn logs a warning message with optional fields
	Warn(msg string, fields ...Field)

	// Error logs an error message with optional fields
	Error(msg string, fields ...Field)

	// With returns a new logger with the given fields attached
	With(fields ...Field) Logger

	// WithContext returns a new logger with context information attached
	WithContext(ctx context.Context) Logger
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// NewField creates a new Field with the given key and value
func NewField(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// Convenience functions for creating common field types
func String(key, value string) Field {
	return NewField(key, value)
}

func Int(key string, value int) Field {
	return NewField(key, value)
}

func Int64(key string, value int64) Field {
	return NewField(key, value)
}

func Float64(key string, value float64) Field {
	return NewField(key, value)
}

func Bool(key string, value bool) Field {
	return NewField(key, value)
}

func Error(err error) Field {
	return NewField("error", err)
}

func Duration(key string, value interface{}) Field {
	return NewField(key, value)
}

