package http

import (
	"context"
	"encoding/json"
	"github.com/SpiridonovDaniil/Distributed-config/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type service interface {
	Create(ctx context.Context, req *domain.Request) error
	Get(ctx context.Context, key string) (json.RawMessage, error)
	GetVersions(ctx context.Context, key string) ([]*domain.Config, error)
	Update(ctx context.Context, req *domain.Request) error
	Delete(ctx context.Context, key string, version int) error
}

func NewServer(service service) *fiber.App {
	f := fiber.New()

	f.Post("/config", createHandler(service))
	f.Get("/config", getHandler(service))
	f.Get("/config/versions", getVersionsHandler(service))
	f.Put("/config", putHandler(service))
	f.Delete("/config", deleteHandler(service))

	return f
}
