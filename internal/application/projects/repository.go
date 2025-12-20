package projects

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mrthoabby/portfolio-api/internal/repository"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *repository.Database) *Repository {
	return &Repository{
		collection: db.Database.Collection("projects"),
	}
}

func (r *Repository) GetByProfileID(ctx context.Context, profileID string) ([]Project, error) {
	filter := bson.M{
		"profileId": profileID,
		"visible":   true,
	}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []Project
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}
