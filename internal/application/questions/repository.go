package questions

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
		store: dataSource.Store("questions"),
	}
}

func (r *Repository) Create(ctx context.Context, profileID, message, ip string) (*Question, error) {
	now := time.Now()
	newQuestion := &Question{
		ID:        uuid.New().String(),
		ProfileID: profileID,
		Message:   message,
		IP:        ip,
		CreatedAt: now,
	}

	err := r.store.InsertOne(ctx, newQuestion)
	if err != nil {
		return nil, err
	}

	return newQuestion, nil
}

