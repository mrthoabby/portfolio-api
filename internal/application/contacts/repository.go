package contacts

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mrthoabby/portfolio-api/internal/common/contracts"
)

type Repository struct {
	store contracts.Store
}

func NewRepository(dataSource contracts.DataSource) *Repository {
	return &Repository{
		store: dataSource.Store("contacts"),
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
		Contacted: false,
		CreatedAt: now,
	}

	err := r.store.InsertOne(ctx, newContact)
	if err != nil {
		return nil, err
	}

	return newContact, nil
}
