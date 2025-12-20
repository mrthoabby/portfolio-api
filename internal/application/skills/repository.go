package skills

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
		collection: db.Database.Collection("skills"),
	}
}

func (r *Repository) GetByProfileID(ctx context.Context, profileID string) ([]Skill, error) {
	filter := bson.M{"profileId": profileID}
	opts := options.Find().SetSort(bson.D{
		{Key: "category", Value: 1},
		{Key: "name", Value: 1},
		{Key: "proficiency", Value: 1},
	})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var skills []Skill
	if err := cursor.All(ctx, &skills); err != nil {
		return nil, err
	}

	return skills, nil
}
