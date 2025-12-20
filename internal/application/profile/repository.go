package profile

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mrthoabby/portfolio-api/internal/repository"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *repository.Database) *Repository {
	return &Repository{
		collection: db.Database.Collection("profiles"),
	}
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Profile, error) {
	var profile Profile
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return &profile, nil
}

func (r *Repository) Exists(ctx context.Context, id string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
