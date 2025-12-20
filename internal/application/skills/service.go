package skills

import (
	"context"
	"errors"

	"github.com/mrthoabby/portfolio-api/internal/application/profile"
)

type Service struct {
	repo           *Repository
	profileService *profile.Service
}

func NewService(repo *Repository, profileService *profile.Service) *Service {
	return &Service{
		repo:           repo,
		profileService: profileService,
	}
}

func (s *Service) GetByProfileID(ctx context.Context, profileID string) ([]Skill, error) {
	// Verify profile exists
	exists, err := s.profileService.Exists(ctx, profileID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("profile not found")
	}

	skills, err := s.repo.GetByProfileID(ctx, profileID)
	if err != nil {
		return nil, err
	}
	return skills, nil
}

