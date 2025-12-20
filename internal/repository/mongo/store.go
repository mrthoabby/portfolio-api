package mongo

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mrthoabby/portfolio-api/internal/common/types"
)

// Store implements contract.Store for MongoDB collections.
type Store struct {
	collection *mongo.Collection
}

// NewStore creates a new Store wrapper for a MongoDB collection.
func NewStore(collection *mongo.Collection) *Store {
	return &Store{collection: collection}
}

// FindOne finds a single record matching the filter and decodes it into result.
func (s *Store) FindOne(ctx context.Context, filter map[string]interface{}, result interface{}) error {
	bsonFilter := toBsonM(filter)
	err := s.collection.FindOne(ctx, bsonFilter).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return types.ErrNotFound{Message: "record not found"}
		}
		return err
	}
	return nil
}

// FindMany finds all records matching the filter with optional sorting.
func (s *Store) FindMany(ctx context.Context, filter map[string]interface{}, sortFields []string, results interface{}) error {
	bsonFilter := toBsonM(filter)
	opts := options.Find()

	if len(sortFields) > 0 {
		sortDoc := bson.D{}
		for _, field := range sortFields {
			order := 1
			if strings.HasPrefix(field, "-") {
				order = -1
				field = strings.TrimPrefix(field, "-")
			}
			sortDoc = append(sortDoc, bson.E{Key: field, Value: order})
		}
		opts.SetSort(sortDoc)
	}

	cursor, err := s.collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, results)
}

// InsertOne inserts a single record into the store.
func (s *Store) InsertOne(ctx context.Context, record interface{}) error {
	_, err := s.collection.InsertOne(ctx, record)
	return err
}

// CountRecords counts the number of records matching the filter.
func (s *Store) CountRecords(ctx context.Context, filter map[string]interface{}) (int64, error) {
	bsonFilter := toBsonM(filter)
	return s.collection.CountDocuments(ctx, bsonFilter)
}

// UpdateOne updates a single record matching the filter.
func (s *Store) UpdateOne(ctx context.Context, filter map[string]interface{}, update map[string]interface{}) error {
	bsonFilter := toBsonM(filter)
	bsonUpdate := bson.M{"$set": toBsonM(update)}
	_, err := s.collection.UpdateOne(ctx, bsonFilter, bsonUpdate)
	return err
}

// toBsonM converts a map[string]interface{} to bson.M.
func toBsonM(m map[string]interface{}) bson.M {
	if m == nil {
		return bson.M{}
	}
	return bson.M(m)
}
