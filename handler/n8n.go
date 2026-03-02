package handler

import (
	"github.com/gofiber/fiber/v2"
)

// CreateWorkflowRequest represents the request to create a new workflow
type CreateWorkflowRequest struct {
	Name     string `json:"name"`
	Prompt   string `json:"prompt"`
	BotToken string `json:"token"`
}

// UpdateWorkflowRequest represents the request to update an existing workflow
type UpdateWorkflowRequest struct {
	Prompt string `json:"prompt"` // Natural language instruction for modification
}

// WorkflowResponse represents a workflow response
type WorkflowResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	N8NID       string `json:"n8n_id,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateWorkflow handles POST /api/n8n/workflows
func (h *Handler) CreateWorkflow(c *fiber.Ctx) error {
	userID, ok := c.Locals("id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	var req CreateWorkflowRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Name == "" || req.Prompt == "" || req.BotToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name, prompt and bot_token are required",
		})
	}

	workflow, err := h.service.CreateWorkflowWithAI(userID, req.Name, req.Prompt, req.BotToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to create workflow",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(WorkflowResponse{
		ID:          workflow.ID,
		Name:        workflow.Name,
		Description: workflow.Description,
		Active:      workflow.Active,
		N8NID:       workflow.WorkflowID,
		CreatedAt:   workflow.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   workflow.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// UpdateWorkflow handles PUT /api/n8n/workflows/:id
func (h *Handler) UpdateWorkflow(c *fiber.Ctx) error {
	userID, ok := c.Locals("id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	workflowID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid workflow ID",
		})
	}

	// Verify ownership
	workflow, err := h.service.GetWorkflowByID(uint(workflowID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "workflow not found",
		})
	}

	if workflow.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "access denied",
		})
	}

	var req UpdateWorkflowRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "prompt is required",
		})
	}

	updatedWorkflow, err := h.service.UpdateWorkflowWithAI(uint(workflowID), req.Prompt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to update workflow",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(WorkflowResponse{
		ID:          updatedWorkflow.ID,
		Name:        updatedWorkflow.Name,
		Description: updatedWorkflow.Description,
		Active:      updatedWorkflow.Active,
		N8NID:       updatedWorkflow.WorkflowID,
		CreatedAt:   updatedWorkflow.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   updatedWorkflow.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// GetWorkflows handles GET /api/n8n/workflows
func (h *Handler) GetWorkflows(c *fiber.Ctx) error {
	userID, ok := c.Locals("id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	workflows, err := h.service.GetUserWorkflows(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to get workflows",
			"details": err.Error(),
		})
	}

	response := make([]WorkflowResponse, len(workflows))
	for i, wf := range workflows {
		response[i] = WorkflowResponse{
			ID:          wf.ID,
			Name:        wf.Name,
			Description: wf.Description,
			Active:      wf.Active,
			N8NID:       wf.WorkflowID,
			CreatedAt:   wf.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   wf.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetWorkflow handles GET /api/n8n/workflows/:id
func (h *Handler) GetWorkflow(c *fiber.Ctx) error {
	userID, ok := c.Locals("id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	workflowID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid workflow ID",
		})
	}

	workflow, err := h.service.GetWorkflowByID(uint(workflowID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "workflow not found",
		})
	}

	if workflow.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "access denied",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":            workflow.ID,
		"name":          workflow.Name,
		"description":   workflow.Description,
		"workflow_json": workflow.WorkflowJSON,
		"chat_history":  workflow.ChatHistory,
		"active":        workflow.Active,
		"n8n_id":        workflow.WorkflowID,
		"created_at":    workflow.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"updated_at":    workflow.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// DeleteWorkflow handles DELETE /api/n8n/workflows/:id
func (h *Handler) DeleteWorkflow(c *fiber.Ctx) error {
	userID, ok := c.Locals("id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not authenticated",
		})
	}

	workflowID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid workflow ID",
		})
	}

	// Verify ownership
	workflow, err := h.service.GetWorkflowByID(uint(workflowID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "workflow not found",
		})
	}

	if workflow.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "access denied",
		})
	}

	if err := h.service.DeleteWorkflow(uint(workflowID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to delete workflow",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "workflow deleted successfully",
	})
}

// RegisterWebhook handles POST /api/n8n/workflows/:id/webhook
func (h *Handler) RegisterWebhook(c *fiber.Ctx) error {
	userID, ok := c.Locals("id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not authenticated"})
	}

	workflowID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid workflow ID"})
	}

	workflow, err := h.service.GetWorkflowByID(uint(workflowID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "workflow not found"})
	}

	if workflow.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "access denied"})
	}

	var req struct {
		BotToken string `json:"token"`
	}
	if err := c.BodyParser(&req); err != nil || req.BotToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token is required"})
	}

	webhookURL, err := h.service.RegisterWorkflowWebhook(uint(workflowID), req.BotToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to register webhook",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"webhook_url": webhookURL})
}
