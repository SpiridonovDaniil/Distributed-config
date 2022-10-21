package http

import (
	"context"
	"encoding/json"

	"github.com/SpiridonovDaniil/Distributed-config/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type service interface {
	Create(ctx context.Context, req *domain.Config) error
	Get(ctx context.Context, key string) (json.RawMessage, error)
}

func NewServer(service service) *fiber.App {
	f := fiber.New()

	f.Post("/config", createHandler(service))
	f.Get("/config", getHandler(service))
	f.Put("/config")
	f.Delete("/config")

	return f
}
