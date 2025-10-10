package main

import (
	"github.com/gofiber/fiber/v2"

	"main/handler"
)

func SetupRoutes(app *fiber.App, h *handler.Handler) {
	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
