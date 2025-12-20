package contracts

import "context"

// Store defines the contract for data storage operations.
// This abstraction allows switching between different storage implementations
// (MongoDB collections, PostgreSQL tables, files, in-memory, etc.)
// without changing the application layer.
type Store interface {
	// FindOne finds a single record matching the filter and decodes it into result.
	// Returns ErrNotFound if no record matches.
	FindOne(ctx context.Context, filter map[string]interface{}, result interface{}) error

	// FindMany finds all records matching the filter with optional sorting.
	// sortFields is a slice of field names; prefix with "-" for descending order.
	// Example: []string{"category", "-createdAt"} sorts by category ASC, then createdAt DESC.
	FindMany(ctx context.Context, filter map[string]interface{}, sortFields []string, results interface{}) error

	// InsertOne inserts a single record into the store.
	InsertOne(ctx context.Context, record interface{}) error

	// CountRecords counts the number of records matching the filter.
	CountRecords(ctx context.Context, filter map[string]interface{}) (int64, error)

	// UpdateOne updates a single record matching the filter.
	UpdateOne(ctx context.Context, filter map[string]interface{}, update map[string]interface{}) error
}
