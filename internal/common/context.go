package common

import "context"

// ContextKey types for type-safe context values
type clientIPKey struct{}

// WithClientIP stores the client IP in the context
func WithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, clientIPKey{}, ip)
}

// ClientIPFromContext retrieves the client IP from context
func ClientIPFromContext(ctx context.Context) string {
	if ip, ok := ctx.Value(clientIPKey{}).(string); ok {
		return ip
	}
	return ""
}

