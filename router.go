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

	// API routes
	api := app.Group("/api")

	// Agent routes
	agent := api.Group("/agent")
	agent.Post("/create", h.FirebaseAuth, h.CreateBot)      // POST /api/agent/create (protected with JWT)
	agent.Get("/list", h.FirebaseAuth, h.GetUserBots)       // GET /api/agent/list (protected with JWT)
	agent.Delete("/:agentId", h.FirebaseAuth, h.DeleteBot)  // DELETE /api/agent/:agentId (protected with JWT)
	agent.Post("/chat", h.Chat)                             // POST /api/agent/chat (public for widget)
}
