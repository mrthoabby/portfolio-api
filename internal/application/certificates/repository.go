package certificates

import (
	"context"

	"github.com/mrthoabby/portfolio-api/internal/common/contracts"
)

type Repository struct {
	store contracts.Store
}

func NewRepository(dataSource contracts.DataSource) *Repository {
	return &Repository{
		store: dataSource.Store("certificates"),
	}
}

func (r *Repository) GetByProfileID(ctx context.Context, profileID string) ([]Certificate, error) {
	var certificates []Certificate
	filter := map[string]interface{}{"profileId": profileID}
	sortFields := []string{"name"}

	err := r.store.FindMany(ctx, filter, sortFields, &certificates)
	if err != nil {
		return nil, err
	}

	return certificates, nil
}
