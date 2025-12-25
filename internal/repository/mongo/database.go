package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mrthoabby/portfolio-api/internal/common/contracts"
	"github.com/mrthoabby/portfolio-api/internal/common/logger"
)

const (
	// Connection timeouts
	connectionTimeout = 10 * time.Second
	pingTimeout       = 30 * time.Second
	disconnectTimeout = 5 * time.Second

	// Connection pool settings
	maxPoolSize     = 100
	minPoolSize     = 10
	maxConnIdleTime = 30 * time.Second
)

// DataSource implements contract.DataSource for MongoDB.
type DataSource struct {
	client   *mongo.Client
	database *mongo.Database
}

// Ensure DataSource implements contract.DataSource.
var _ contracts.DataSource = (*DataSource)(nil)

// NewDataSource creates a new MongoDB data source connection using the provided logger.
// The logger is required and cannot be nil.
func NewDataSource(connectionString, databaseName string, logs logger.Logger) (*DataSource, error) {
	sanitizedURL := sanitizeMongoURLForLogging(connectionString)
	logs.Info("Attempting to connect to MongoDB",
		logger.String("url", sanitizedURL),
		logger.String("database", databaseName),
	)

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(connectionString)
	logs.Debug("Connection options configured")

	clientOptions.SetMaxPoolSize(maxPoolSize)
	clientOptions.SetMinPoolSize(minPoolSize)
	clientOptions.SetMaxConnIdleTime(maxConnIdleTime)

	logs.Debug("Attempting connection",
		logger.String("timeout", connectionTimeout.String()),
	)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logs.Error("Connection failed", logger.Error(err))
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	logs.Debug("Connection established, verifying with ping")
	pingCtx, pingCancel := context.WithTimeout(context.Background(), pingTimeout)
	defer pingCancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		logs.Error("Ping failed", logger.Error(err))
		client.Disconnect(context.Background())
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logs.Info("Successfully connected and verified to database",
		logger.String("database", databaseName),
	)

	db := client.Database(databaseName)

	return &DataSource{
		client:   client,
		database: db,
	}, nil
}

// Store returns a Store for the given name (MongoDB collection).
func (d *DataSource) Store(name string) contracts.Store {
	return NewStore(d.database.Collection(name))
}

// Close closes the data source connection.
func (d *DataSource) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), disconnectTimeout)
	defer cancel()
	return d.client.Disconnect(ctx)
}

// Ping verifies the data source connection is alive.
func (d *DataSource) Ping(ctx context.Context) error {
	return d.client.Ping(ctx, nil)
}
