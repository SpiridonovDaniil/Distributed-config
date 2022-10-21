package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SpiridonovDaniil/Distributed-config/internal/domain"
	"github.com/SpiridonovDaniil/Distributed-config/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, config *domain.Config) error {
	err := s.repo.Create(ctx, config.Service, config.Data)
	if err != nil {
		return fmt.Errorf("[create] failed to create config, error: %w", err)
	}

	return nil
}

func (s *Service) Get(ctx context.Context, key string) (json.RawMessage, error) {
	resp, err := s.repo.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("[get] failed to get config, error: %w", err)
	}

	return resp, nil
}
