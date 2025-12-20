package skills

import (
	"context"

	"github.com/mrthoabby/portfolio-api/internal/common/contracts"
)

type Repository struct {
	store contracts.Store
}

func NewRepository(dataSource contracts.DataSource) *Repository {
	return &Repository{
		store: dataSource.Store("skills"),
	}
}

func (r *Repository) GetByProfileID(ctx context.Context, profileID string) ([]Skill, error) {
	var skills []Skill
	filter := map[string]interface{}{"profileId": profileID}
	sortFields := []string{"category", "name", "proficiency"}

	err := r.store.FindMany(ctx, filter, sortFields, &skills)
	if err != nil {
		return nil, err
	}

	return skills, nil
}
