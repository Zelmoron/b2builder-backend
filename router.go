package main

import (
	"github.com/gofiber/fiber/v2"

	"main/handler"
)

func SetupRoutes(app *fiber.App, h *handler.Handler) {

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/reg", h.RegisterUser)
	//api := app.Group("/api")
}
