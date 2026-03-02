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
	app.Post("/:id/webhook", h.RegisterWebhook)

	// API routes
	api := app.Group("/api")

	// Agent routes
	agent := api.Group("/agent")
	agent.Post("/create", h.FirebaseAuth, h.CreateBot)     // POST /api/agent/create (protected with JWT)
	agent.Get("/list", h.FirebaseAuth, h.GetUserBots)      // GET /api/agent/list (protected with JWT)
	agent.Delete("/:agentId", h.FirebaseAuth, h.DeleteBot) // DELETE /api/agent/:agentId (protected with JWT)
	agent.Post("/chat", h.Chat)                            // POST /api/agent/chat (public for widget)

	// N8N AI Workflow routes
	workflows := api.Group("/n8n/workflows", h.FirebaseAuth)
	workflows.Post("/", h.CreateWorkflow)      // POST /api/n8n/workflows - Create new workflow with AI
	workflows.Get("/", h.GetWorkflows)         // GET /api/n8n/workflows - Get all user workflows
	workflows.Get("/:id", h.GetWorkflow)       // GET /api/n8n/workflows/:id - Get specific workflow with chat history
	workflows.Put("/:id", h.UpdateWorkflow)    // PUT /api/n8n/workflows/:id - Update workflow with AI
	workflows.Delete("/:id", h.DeleteWorkflow) // DELETE /api/n8n/workflows/:id - Delete workflow
	// POST /api/n8n/workflows/:id/webhook - Force register Telegram webhook
}
