package repository

import (
	"context"
	"encoding/json"
)

type Repository interface {
	Create(ctx context.Context, key string, metaData json.RawMessage) error
	Get(ctx context.Context, key string) (json.RawMessage, error)
	Update(ctx context.Context, key string, metaData json.RawMessage) error
	Delete(ctx context.Context, key string) error
}
