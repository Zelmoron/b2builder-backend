package services

import (
	"encoding/json"
	"fmt"
	"main/models"
	"time"
)

// CreateWorkflowWithAI creates a new workflow using AI based on user's natural language prompt
func (s *Service) CreateWorkflowWithAI(userID uint, workflowName, userPrompt string) (*models.N8NWorkflow, error) {
	// Build messages for AI
	messages := []OpenRouterMessage{
		{
			Role:    "system",
			Content: N8N_WORKFLOW_SYSTEM_PROMPT,
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	}

	// Call AI to generate workflow JSON
	aiResponse, err := s.callOpenRouter(messages)
	if err != nil {
		return nil, fmt.Errorf("AI generation failed: %w", err)
	}

	// Extract JSON from markdown code blocks if present
	cleanJSON := extractJSON(aiResponse)

	// Parse AI response as N8N workflow
	var workflowData map[string]interface{}
	if err := json.Unmarshal([]byte(cleanJSON), &workflowData); err != nil {
		return nil, fmt.Errorf("failed to parse AI response as JSON: %w (response: %s)", err, cleanJSON)
	}

	// Create chat history
	chatHistory := []models.WorkflowChatMessage{
		{
			Role:      "user",
			Content:   userPrompt,
			Timestamp: time.Now(),
		},
		{
			Role:      "assistant",
			Content:   aiResponse,
			Timestamp: time.Now(),
		},
	}

	chatHistoryJSON, err := json.Marshal(chatHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chat history: %w", err)
	}

	workflowJSON, err := json.Marshal(workflowData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow data: %w", err)
	}

	// Create workflow in database
	dbWorkflow := &models.N8NWorkflow{
		UserID:       userID,
		Name:         workflowName,
		Description:  userPrompt,
		WorkflowJSON: workflowJSON,
		ChatHistory:  chatHistoryJSON,
		Active:       false,
	}

	if err := s.repo.CreateN8NWorkflow(dbWorkflow); err != nil {
		return nil, fmt.Errorf("failed to save workflow to database: %w", err)
	}

	// Try to create in N8N
	n8nWorkflow := s.convertMapToN8NWorkflow(workflowData)
	n8nResp, err := s.createN8NWorkflow(n8nWorkflow)
	if err == nil && n8nResp != nil {
		// Update with N8N workflow ID
		dbWorkflow.WorkflowID = n8nResp.ID
		dbWorkflow.Active = n8nResp.Active
		s.repo.UpdateN8NWorkflow(dbWorkflow)
	}

	return dbWorkflow, nil
}

// UpdateWorkflowWithAI updates an existing workflow using AI
func (s *Service) UpdateWorkflowWithAI(workflowID uint, userPrompt string) (*models.N8NWorkflow, error) {
	// Get existing workflow
	dbWorkflow, err := s.repo.GetN8NWorkflowByID(workflowID)
	if err != nil {
		return nil, fmt.Errorf("workflow not found: %w", err)
	}

	// Parse existing chat history
	var chatHistory []models.WorkflowChatMessage
	if err := json.Unmarshal(dbWorkflow.ChatHistory, &chatHistory); err != nil {
		return nil, fmt.Errorf("failed to parse chat history: %w", err)
	}

	// Build messages for AI (include history + current workflow)
	messages := []OpenRouterMessage{
		{
			Role:    "system",
			Content: N8N_WORKFLOW_SYSTEM_PROMPT,
		},
	}

	// Add chat history
	for _, msg := range chatHistory {
		messages = append(messages, OpenRouterMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Add current workflow context and new request
	contextPrompt := fmt.Sprintf("Current workflow JSON:\n%s\n\nUser request: %s", string(dbWorkflow.WorkflowJSON), userPrompt)
	messages = append(messages, OpenRouterMessage{
		Role:    "user",
		Content: contextPrompt,
	})

	// Call AI to update workflow
	aiResponse, err := s.callOpenRouter(messages)
	if err != nil {
		return nil, fmt.Errorf("AI generation failed: %w", err)
	}

	// Extract JSON from markdown code blocks if present
	cleanJSON := extractJSON(aiResponse)

	// Parse AI response
	var workflowData map[string]interface{}
	if err := json.Unmarshal([]byte(cleanJSON), &workflowData); err != nil {
		return nil, fmt.Errorf("failed to parse AI response as JSON: %w (response: %s)", err, cleanJSON)
	}

	// Update chat history
	chatHistory = append(chatHistory,
		models.WorkflowChatMessage{
			Role:      "user",
			Content:   userPrompt,
			Timestamp: time.Now(),
		},
		models.WorkflowChatMessage{
			Role:      "assistant",
			Content:   aiResponse,
			Timestamp: time.Now(),
		},
	)

	chatHistoryJSON, err := json.Marshal(chatHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chat history: %w", err)
	}

	workflowJSON, err := json.Marshal(workflowData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow data: %w", err)
	}

	// Update database
	dbWorkflow.WorkflowJSON = workflowJSON
	dbWorkflow.ChatHistory = chatHistoryJSON
	dbWorkflow.UpdatedAt = time.Now()

	if err := s.repo.UpdateN8NWorkflow(dbWorkflow); err != nil {
		return nil, fmt.Errorf("failed to update workflow in database: %w", err)
	}

	// Try to update in N8N if workflow ID exists
	if dbWorkflow.WorkflowID != "" {
		n8nWorkflow := s.convertMapToN8NWorkflow(workflowData)
		n8nResp, err := s.updateN8NWorkflowByID(dbWorkflow.WorkflowID, n8nWorkflow)
		if err == nil && n8nResp != nil {
			dbWorkflow.Active = n8nResp.Active
			s.repo.UpdateN8NWorkflow(dbWorkflow)
		}
	}

	return dbWorkflow, nil
}

// GetUserWorkflows returns all workflows for a user
func (s *Service) GetUserWorkflows(userID uint) ([]models.N8NWorkflow, error) {
	return s.repo.GetN8NWorkflowsByUserID(userID)
}

// GetWorkflowByID returns a specific workflow
func (s *Service) GetWorkflowByID(workflowID uint) (*models.N8NWorkflow, error) {
	return s.repo.GetN8NWorkflowByID(workflowID)
}

// DeleteWorkflow deletes a workflow
func (s *Service) DeleteWorkflow(workflowID uint) error {
	// Get workflow to check if it has N8N ID
	workflow, err := s.repo.GetN8NWorkflowByID(workflowID)
	if err == nil && workflow.WorkflowID != "" {
		// Try to delete from N8N (ignore errors)
		s.deleteN8NWorkflowByID(workflow.WorkflowID)
	}

	return s.repo.DeleteN8NWorkflow(workflowID)
}

// Helper function to convert map to N8NWorkflowRequest
func (s *Service) convertMapToN8NWorkflow(data map[string]interface{}) *N8NWorkflowRequest {
	workflow := &N8NWorkflowRequest{
		Name:        "",
		Nodes:       []interface{}{},
		Connections: map[string]interface{}{},
		Settings:    map[string]interface{}{},
	}

	if name, ok := data["name"].(string); ok {
		workflow.Name = name
	}

	if nodes, ok := data["nodes"].([]interface{}); ok {
		workflow.Nodes = nodes
	}

	if connections, ok := data["connections"].(map[string]interface{}); ok {
		workflow.Connections = connections
	}

	if settings, ok := data["settings"].(map[string]interface{}); ok {
		workflow.Settings = settings
	}

	return workflow
}

// Delete workflow from N8N
func (s *Service) deleteN8NWorkflowByID(workflowID string) error {
	// Implementation for deleting from N8N API (if needed)
	// For now, just a placeholder
	return nil
}