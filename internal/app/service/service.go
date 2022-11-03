package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SpiridonovDaniil/Distributed-config/internal/helpers"

	"github.com/SpiridonovDaniil/Distributed-config/internal/domain"
	"github.com/SpiridonovDaniil/Distributed-config/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, config *domain.Request) error {
	rawData, err := helpers.Converter(config)
	if err != nil {
		return fmt.Errorf("[create] %w", err)
	}

	err = s.repo.Create(ctx, config.Service, rawData)
	if err != nil {
		return fmt.Errorf("[create] failed to create config, error: %w", err)
	}

	return nil
}

func (s *Service) RollBack(ctx context.Context, key string, version int) error {
	err := s.repo.RollBack(ctx, key, version)
	if err != nil {
		return fmt.Errorf("[rollback] failed to rollback config, error: %w", err)
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

func (s *Service) GetVersions(ctx context.Context, key string) ([]*domain.Config, error) {
	resp, err := s.repo.GetVersions(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("[getVersions] failed to get versions config, error: %w", err)
	}

	return resp, nil
}

func (s *Service) Update(ctx context.Context, config *domain.Request) error {
	rawData, err := helpers.Converter(config)
	if err != nil {
		return fmt.Errorf("[update] %w", err)
	}

	err = s.repo.Update(ctx, config.Service, rawData)
	if err != nil {
		return fmt.Errorf("[update] failed to update config, error: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, key string, version int) error {
	err := s.repo.Delete(ctx, key, version)
	if err != nil {
		return fmt.Errorf("[delete] failed to delete config, error: %w", err)
	}

	return nil
}
