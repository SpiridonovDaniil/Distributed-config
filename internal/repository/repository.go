package repository

import (
	"context"
	"encoding/json"
	"github.com/SpiridonovDaniil/Distributed-config/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, key string, metaData json.RawMessage) error
	RollBack(ctx context.Context, key string, version int) error
	Get(ctx context.Context, key string) (json.RawMessage, error)
	GetVersions(ctx context.Context, key string) ([]*domain.Config, error)
	Update(ctx context.Context, key string, metaData json.RawMessage) error
	Delete(ctx context.Context, key string, version int) error
}
