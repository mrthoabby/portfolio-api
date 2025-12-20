package profile

import (
	"context"
	"errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(ctx context.Context, id string) (*Profile, error) {
	profile, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("profile not found")
	}
	return profile, nil
}

func (s *Service) Exists(ctx context.Context, id string) (bool, error) {
	return s.repo.Exists(ctx, id)
}
