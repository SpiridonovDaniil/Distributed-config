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
	var metadata []map[string]interface{}
	err := json.Unmarshal(config.Data, &metadata)
	if err != nil {
		// todo asdasda
	}

	resultMeta := make(map[string]interface{})
	for _, data := range metadata {
		for key, val := range data {
			resultMeta[key] = val
		}
	}

	rawData, _ := json.Marshal(resultMeta)

	err = s.repo.Create(ctx, config.Service, rawData)
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

func (s *Service) Update(ctx context.Context, config *domain.Config) error {
	err := s.repo.Update(ctx, config.Service, config.Data)
	if err != nil {
		return fmt.Errorf("[update] failed to update config, error: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, key string) error {
	err := s.repo.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("[delete] failed to delete config, error: %w", err)
	}

	return nil
}
