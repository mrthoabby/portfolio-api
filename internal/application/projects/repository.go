package projects

import (
	"context"

	"github.com/mrthoabby/portfolio-api/internal/common/contracts"
)

type Repository struct {
	store contracts.Store
}

func NewRepository(dataSource contracts.DataSource) *Repository {
	return &Repository{
		store: dataSource.Store("projects"),
	}
}

func (r *Repository) GetByProfileID(ctx context.Context, profileID string) ([]Project, error) {
	var projects []Project
	filter := map[string]interface{}{
		"profileId": profileID,
		"visible":   true,
	}
	sortFields := []string{"-createdAt"} // "-" prefix for descending order

	err := r.store.FindMany(ctx, filter, sortFields, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}
