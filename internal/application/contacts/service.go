package contacts

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

func (s *Service) Create(ctx context.Context, profileID string, contact *Request) (*Contact, error) {
	// Verify profile exists
	exists, err := s.profileService.Exists(ctx, profileID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("profile not found")
	}

	createdContact, err := s.repo.Create(ctx, profileID, contact)
	if err != nil {
		return nil, err
	}
	return createdContact, nil
}
