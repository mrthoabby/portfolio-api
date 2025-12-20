package profile

import (
	"context"

	"github.com/mrthoabby/portfolio-api/internal/common/contracts"
	"github.com/mrthoabby/portfolio-api/internal/common/types"
)

type Repository struct {
	store contracts.Store
}

func NewRepository(dataSource contracts.DataSource) *Repository {
	return &Repository{
		store: dataSource.Store("profiles"),
	}
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Profile, error) {
	var profile Profile
	err := r.store.FindOne(ctx, map[string]interface{}{"_id": id}, &profile)
	if err != nil {
		if types.IsNotFoundError(err) {
			return nil, types.ErrNotFound{Message: "profile not found"}
		}
		return nil, err
	}
	return &profile, nil
}

func (r *Repository) Exists(ctx context.Context, id string) (bool, error) {
	count, err := r.store.CountRecords(ctx, map[string]interface{}{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
