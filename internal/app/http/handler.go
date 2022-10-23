package http

import (
	"fmt"
	"net/http"

	"github.com/SpiridonovDaniil/Distributed-config/internal/domain"
	"github.com/gofiber/fiber/v2"
)

func createHandler(service service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var req domain.Config
		err := ctx.BodyParser(&req)
		if err != nil {
			return fmt.Errorf("[createHandler] failed to parse request, error: %w", err)
		}
		err = service.Create(ctx.Context(), &req)
		if err != nil {
			return fmt.Errorf("[createHandler] %w", err)
		}
		ctx.Status(http.StatusCreated)

		return nil
	}
}

func getHandler(service service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		key := ctx.Get("service")
		resp, err := service.Get(ctx.Context(), key)
		if err != nil {
			return fmt.Errorf("[getHandler] %w", err)
		}
		err = ctx.JSON(resp)
		if err != nil {
			return fmt.Errorf("[getHandler] failed to return JSON answer, error: %w", err)
		}

		return nil
	}
}

func putHandler(service service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var req domain.Config
		err := ctx.BodyParser(&req)
		if err != nil {
			return fmt.Errorf("[putHandler] failed to parse request, error: %w", err)
		}
		err = service.Update(ctx.Context(), &req)
		if err != nil {
			return fmt.Errorf("[putHandler] %w", err)
		}
		ctx.Status(http.StatusOK)

		return nil
	}
}

func deleteHandler(service service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		key := ctx.Get("service")
		err := service.Delete(ctx.Context(), key)
		if err != nil {
			return fmt.Errorf("[deleteHandler] %w", err)
		}
		ctx.Status(http.StatusOK)

		return nil
	}
}