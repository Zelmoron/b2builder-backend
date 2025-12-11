package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"main/models"
)

// CreateBotRequest represents the request body for creating a bot
type CreateBotRequest struct {
	Name               string           `json:"name"`
	Type               string           `json:"type"`
	ProductDescription string           `json:"productDescription"`
	FAQ                []models.FAQItem `json:"faq"`
}

// CreateBotResponse represents the response for bot creation
type CreateBotResponse struct {
	BotID   string `json:"botId"`
	AgentID uint   `json:"agentId"`
	Status  string `json:"status"`
}

// ChatRequest represents the request body for chat messages
type ChatRequest struct {
	BotID     string `json:"botId"`
	SessionID string `json:"sessionId"`
	Message   string `json:"message"`
}

// ChatResponse represents the response for chat messages
type ChatResponse struct {
	Reply string `json:"reply"`
}

// AgentListItem represents a single agent in the list response
type AgentListItem struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	BotID     string `json:"bot_id"`
	CreatedAt string `json:"created_at"`
}

// GetUserBotsResponse represents the response for listing user's agents
type GetUserBotsResponse struct {
	Agents []AgentListItem `json:"agents"`
}

// CreateBot handles POST /api/agent/create
func (h *Handler) CreateBot(c *fiber.Ctx) error {
	// Get user ID from locals (set by FirebaseAuth middleware)
	userID, ok := c.Locals("fbUID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Parse request body
	var req CreateBotRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name is required",
		})
	}
	if req.Type == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "type is required",
		})
	}

	// Create bot
	bot, err := h.service.CreateBot(userID, req.Name, req.Type, req.ProductDescription, req.FAQ)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create bot",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(CreateBotResponse{
		BotID:   bot.BotID,
		AgentID: bot.ID,
		Status:  "success",
	})
}

// Chat handles POST /api/agent/chat
func (h *Handler) Chat(c *fiber.Ctx) error {
	// Parse request body
	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate required fields
	if req.BotID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "botId is required",
		})
	}
	if req.SessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "sessionId is required",
		})
	}
	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "message is required",
		})
	}

	// Process chat message
	reply, err := h.service.ProcessChatMessage(req.BotID, req.SessionID, req.Message)
	if err != nil {
		if err.Error() == "bot not found or inactive" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "bot not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to process message",
		})
	}

	// Return response
	return c.JSON(ChatResponse{
		Reply: reply,
	})
}

// GetUserBots handles GET /api/agent/list (endpoint for listing user's bots)
func (h *Handler) GetUserBots(c *fiber.Ctx) error {
	// Get user ID from locals (set by FirebaseAuth middleware)
	userID, ok := c.Locals("fbUID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get user's bots
	bots, err := h.service.GetUserBots(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve bots",
		})
	}

	// Convert to response format
	agents := make([]AgentListItem, len(bots))
	for i, bot := range bots {
		agents[i] = AgentListItem{
			ID:        bot.ID,
			Name:      bot.Name,
			Type:      bot.Type,
			BotID:     bot.BotID,
			CreatedAt: bot.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), // ISO 8601 format
		}
	}

	return c.JSON(GetUserBotsResponse{
		Agents: agents,
	})
}

// DeleteBot handles DELETE /api/agent/:agentId (endpoint for deleting a bot)
func (h *Handler) DeleteBot(c *fiber.Ctx) error {
	// Get user ID from locals (set by FirebaseAuth middleware)
	userID, ok := c.Locals("fbUID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get agentId from URL parameters
	agentIDStr := c.Params("agentId")
	if agentIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "agentId is required",
		})
	}

	// Convert agentId to uint
	agentID64, err := strconv.ParseUint(agentIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid agentId",
		})
	}
	agentID := uint(agentID64)

	// Delete bot
	if err := h.service.DeleteBot(agentID, userID); err != nil {
		if err.Error() == "bot not found or access denied" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "bot not found or access denied",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete bot",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Агент успешно удален",
	})
}
