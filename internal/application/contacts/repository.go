package contacts

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mrthoabby/portfolio-api/internal/repository"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *repository.Database) *Repository {
	return &Repository{
		collection: db.Database.Collection("contacts"),
	}
}

func (r *Repository) Create(ctx context.Context, profileID string, contact *Request) (*Contact, error) {
	now := time.Now()
	newContact := &Contact{
		ID:        uuid.New().String(),
		ProfileID: profileID,
		Name:      contact.Name,
		Email:     contact.Email,
		Message:   contact.Message,
		CreatedAt: now,
	}

	_, err := r.collection.InsertOne(ctx, newContact)
	if err != nil {
		return nil, err
	}

	return newContact, nil
}
