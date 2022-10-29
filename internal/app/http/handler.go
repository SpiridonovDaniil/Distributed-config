package http

import (
	"fmt"
	"github.com/SpiridonovDaniil/Distributed-config/internal/domain"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

func createHandler(service service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var req domain.Request
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
		key := ctx.Query("service")

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

func getVersionsHandler(service service) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		key := ctx.Query("service")

		resp, err := service.GetVersions(ctx.Context(), key)
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
		var req domain.Request
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
		versionStr := ctx.Query("version")

		var version int
		if versionStr != "" {
			ctx.Status(http.StatusBadRequest)
			return fmt.Errorf("[deleteHandler] version must be positive integer")
		}

		version, err := strconv.Atoi(versionStr)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return fmt.Errorf("[deleteHandler] version must be positive integer")
		}

		err = service.Delete(ctx.Context(), key, version) // todo version
		if err != nil {
			return fmt.Errorf("[deleteHandler] %w", err)
		}
		ctx.Status(http.StatusOK)

		return nil
	}
}
