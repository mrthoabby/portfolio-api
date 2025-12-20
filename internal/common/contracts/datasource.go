package contracts

import "context"

// DataSource defines the contract for database/storage connection management.
// Implementations provide access to named stores (collections, tables, etc.)
// and handle connection lifecycle.
type DataSource interface {
	// Store returns a Store for the given name (collection, table, etc.).
	Store(name string) Store

	// Close closes the data source connection.
	Close() error

	// Ping verifies the data source connection is alive.
	Ping(ctx context.Context) error
}
