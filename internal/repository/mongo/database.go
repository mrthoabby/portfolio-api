package mongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mrthoabby/portfolio-api/internal/common/contracts"
)

// DataSource implements contract.DataSource for MongoDB.
type DataSource struct {
	client   *mongo.Client
	database *mongo.Database
}

// Ensure DataSource implements contract.DataSource.
var _ contracts.DataSource = (*DataSource)(nil)

// NewDataSource creates a new MongoDB data source connection.
func NewDataSource(connectionString, databaseName string) (*DataSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(connectionString)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	caCertPath := "/etc/ssl/certs/ca-certificates.crt"
	if _, err := os.Stat(caCertPath); err == nil {
		caCert, err := os.ReadFile(caCertPath)
		if err == nil {
			caCertPool := x509.NewCertPool()
			if caCertPool.AppendCertsFromPEM(caCert) {
				tlsConfig.RootCAs = caCertPool
			}
		}
	} else {
		if systemCertPool, err := x509.SystemCertPool(); err == nil {
			tlsConfig.RootCAs = systemCertPool
		}
	}

	clientOptions.SetTLSConfig(tlsConfig)

	clientOptions.SetMaxPoolSize(100)
	clientOptions.SetMinPoolSize(10)
	clientOptions.SetMaxConnIdleTime(30 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer pingCancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, err
	}

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.client.Disconnect(ctx)
}

// Ping verifies the data source connection is alive.
func (d *DataSource) Ping(ctx context.Context) error {
	return d.client.Ping(ctx, nil)
}
