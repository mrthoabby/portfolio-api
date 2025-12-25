package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/mrthoabby/portfolio-api/internal/common/scope"
)

type zapLogger struct {
	logger *zap.Logger
}

var _ Logger = (*zapLogger)(nil)

// NewZapLogger creates a new zap-based logger based on the current scope.
// The scope must be initialized with scope.Initialize() before calling this function.
// Production scope uses JSON format, otherwise uses development (text) format
func NewZapLogger() (Logger, error) {
	var zapConfig zap.Config

	if scope.IsProduction() {
		zapConfig = zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapLog, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &zapLogger{logger: zapLog}, nil
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	zapFields := l.convertFields(fields)
	l.logger.Debug(msg, zapFields...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	zapFields := l.convertFields(fields)
	l.logger.Info(msg, zapFields...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	zapFields := l.convertFields(fields)
	l.logger.Warn(msg, zapFields...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	zapFields := l.convertFields(fields)
	l.logger.Error(msg, zapFields...)
}

func (l *zapLogger) With(fields ...Field) Logger {
	zapFields := l.convertFields(fields)
	return &zapLogger{logger: l.logger.With(zapFields...)}
}

// WithContext returns a new logger with context information attached
func (l *zapLogger) WithContext(ctx context.Context) Logger {
	// Extract common context fields if available
	// For now, return the logger as-is
	// This can be extended to extract request ID, user ID, etc. from context
	return l
}

func (l *zapLogger) convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, l.convertField(field))
	}
	return zapFields
}

func (l *zapLogger) convertField(field Field) zap.Field {
	switch v := field.Value.(type) {
	case string:
		return zap.String(field.Key, v)
	case int:
		return zap.Int(field.Key, v)
	case int64:
		return zap.Int64(field.Key, v)
	case float64:
		return zap.Float64(field.Key, v)
	case bool:
		return zap.Bool(field.Key, v)
	case error:
		return zap.Error(v)
	case time.Duration:
		return zap.Duration(field.Key, v)
	case time.Time:
		return zap.Time(field.Key, v)
	default:
		// For any other type, use zap.Any
		return zap.Any(field.Key, v)
	}
}

// Sync flushes any buffered log entries. Should be called before application exit.
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}
