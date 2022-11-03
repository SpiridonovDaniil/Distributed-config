package http

import "github.com/gofiber/fiber/v2"

func HandleErrors(ctx *fiber.Ctx) error {
	err := ctx.Next()

	if err == nil {
		return nil
	}
	
	err = ctx.SendString(err.Error())
	if err != nil {
		return err
	}

	return nil
}
