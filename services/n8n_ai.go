package services

import (
	"encoding/json"
	"fmt"
	"log"
	"main/models"
	"time"

	"github.com/google/uuid"
)

func (s *Service) CreateWorkflowWithAI(userID uint, workflowName, userPrompt, botToken string) (*models.N8NWorkflow, error) {
	log.Printf("[CreateWorkflow] user=%d name=%q prompt=%q", userID, workflowName, userPrompt)

	credID := "bot-cred"
	credName := workflowName + " Bot"

	cred, err := s.createN8NCredential(&N8NCredentialRequest{
		Name: credName,
		Type: "telegramApi",
		Data: map[string]interface{}{
			"accessToken": botToken,
		},
	})
	if err != nil {
		log.Printf("[CreateWorkflow] credential creation failed: %v", err)
	} else {
		credID = cred.ID
		credName = cred.Name
		log.Printf("[CreateWorkflow] credential created: id=%q name=%q", credID, credName)
	}

	systemPrompt := fmt.Sprintf(N8N_WORKFLOW_SYSTEM_PROMPT, credID, credName)

	messages := []OpenRouterMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	log.Printf("[CreateWorkflow] calling AI...")
	aiResponse, err := s.callOpenRouter(messages)
	if err != nil {
		return nil, fmt.Errorf("AI generation failed: %w", err)
	}
	log.Printf("[CreateWorkflow] AI response length=%d", len(aiResponse))
	log.Printf("[CreateWorkflow] AI raw response: %s", aiResponse)

	cleanJSON := extractJSON(aiResponse)
	log.Printf("[CreateWorkflow] extracted JSON length=%d", len(cleanJSON))

	var workflowData map[string]interface{}
	if err := json.Unmarshal([]byte(cleanJSON), &workflowData); err != nil {
		log.Printf("[CreateWorkflow] JSON parse error: %v\nJSON: %s", err, cleanJSON)
		return nil, fmt.Errorf("failed to parse AI response as JSON: %w", err)
	}

	webhookID := injectWebhookIDs(workflowData)
	log.Printf("[CreateWorkflow] webhookId=%s", webhookID)

	chatHistory := []models.WorkflowChatMessage{
		{Role: "user", Content: userPrompt, Timestamp: time.Now()},
		{Role: "assistant", Content: aiResponse, Timestamp: time.Now()},
	}

	chatHistoryJSON, err := json.Marshal(chatHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chat history: %w", err)
	}

	workflowJSON, err := json.Marshal(workflowData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow data: %w", err)
	}

	dbWorkflow := &models.N8NWorkflow{
		UserID:       userID,
		Name:         workflowName,
		Description:  userPrompt,
		WorkflowJSON: workflowJSON,
		ChatHistory:  chatHistoryJSON,
		Active:       true,
	}

	if err := s.repo.CreateN8NWorkflow(dbWorkflow); err != nil {
		return nil, fmt.Errorf("failed to save workflow to database: %w", err)
	}
	log.Printf("[CreateWorkflow] saved to DB id=%d", dbWorkflow.ID)

	n8nWorkflow := s.convertMapToN8NWorkflow(workflowData, workflowName)
	n8nResp, err := s.createN8NWorkflow(n8nWorkflow)
	if err != nil {
		log.Printf("[CreateWorkflow] n8n workflow creation failed: %v", err)
	} else {
		log.Printf("[CreateWorkflow] n8n workflow created: id=%q", n8nResp.ID)
		dbWorkflow.WorkflowID = n8nResp.ID

		if err := s.activateN8NWorkflow(n8nResp.ID); err != nil {
			log.Printf("[CreateWorkflow] n8n workflow activation failed: %v", err)
		} else {
			log.Printf("[CreateWorkflow] n8n workflow activated: id=%q", n8nResp.ID)
			dbWorkflow.Active = true
		}

		s.repo.UpdateN8NWorkflow(dbWorkflow)
	}

	return dbWorkflow, nil
}

func (s *Service) UpdateWorkflowWithAI(workflowID uint, userPrompt string) (*models.N8NWorkflow, error) {
	log.Printf("[UpdateWorkflow] id=%d prompt=%q", workflowID, userPrompt)

	dbWorkflow, err := s.repo.GetN8NWorkflowByID(workflowID)
	if err != nil {
		return nil, fmt.Errorf("workflow not found: %w", err)
	}

	var chatHistory []models.WorkflowChatMessage
	if err := json.Unmarshal(dbWorkflow.ChatHistory, &chatHistory); err != nil {
		return nil, fmt.Errorf("failed to parse chat history: %w", err)
	}

	credID, credName := extractCredentialFromWorkflow(dbWorkflow.WorkflowJSON)
	log.Printf("[UpdateWorkflow] using credential: id=%q name=%q", credID, credName)

	systemPrompt := fmt.Sprintf(N8N_WORKFLOW_SYSTEM_PROMPT, credID, credName)

	messages := []OpenRouterMessage{
		{Role: "system", Content: systemPrompt},
	}

	for _, msg := range chatHistory {
		messages = append(messages, OpenRouterMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	contextPrompt := fmt.Sprintf("Current workflow JSON:\n%s\n\nUser request: %s", string(dbWorkflow.WorkflowJSON), userPrompt)
	messages = append(messages, OpenRouterMessage{Role: "user", Content: contextPrompt})

	log.Printf("[UpdateWorkflow] calling AI...")
	aiResponse, err := s.callOpenRouter(messages)
	if err != nil {
		return nil, fmt.Errorf("AI generation failed: %w", err)
	}
	log.Printf("[UpdateWorkflow] AI response length=%d", len(aiResponse))
	log.Printf("[UpdateWorkflow] AI raw response: %s", aiResponse)

	cleanJSON := extractJSON(aiResponse)
	log.Printf("[UpdateWorkflow] extracted JSON length=%d", len(cleanJSON))

	var workflowData map[string]interface{}
	if err := json.Unmarshal([]byte(cleanJSON), &workflowData); err != nil {
		log.Printf("[UpdateWorkflow] JSON parse error: %v\nJSON: %s", err, cleanJSON)
		return nil, fmt.Errorf("failed to parse AI response as JSON: %w", err)
	}

	injectWebhookIDs(workflowData)

	chatHistory = append(chatHistory,
		models.WorkflowChatMessage{Role: "user", Content: userPrompt, Timestamp: time.Now()},
		models.WorkflowChatMessage{Role: "assistant", Content: aiResponse, Timestamp: time.Now()},
	)

	chatHistoryJSON, err := json.Marshal(chatHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chat history: %w", err)
	}

	workflowJSON, err := json.Marshal(workflowData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow data: %w", err)
	}

	dbWorkflow.WorkflowJSON = workflowJSON
	dbWorkflow.ChatHistory = chatHistoryJSON
	dbWorkflow.UpdatedAt = time.Now()

	if err := s.repo.UpdateN8NWorkflow(dbWorkflow); err != nil {
		return nil, fmt.Errorf("failed to update workflow in database: %w", err)
	}
	log.Printf("[UpdateWorkflow] saved to DB id=%d", dbWorkflow.ID)

	if dbWorkflow.WorkflowID != "" {
		n8nWorkflow := s.convertMapToN8NWorkflow(workflowData, dbWorkflow.Name)
		n8nResp, err := s.updateN8NWorkflowByID(dbWorkflow.WorkflowID, n8nWorkflow)
		if err != nil {
			log.Printf("[UpdateWorkflow] n8n update failed: %v", err)
		} else {
			log.Printf("[UpdateWorkflow] n8n updated: active=%v", n8nResp.Active)
			dbWorkflow.Active = n8nResp.Active
			s.repo.UpdateN8NWorkflow(dbWorkflow)
		}
	}

	return dbWorkflow, nil
}

func (s *Service) RegisterWorkflowWebhook(workflowID uint, botToken string) (string, error) {
	dbWorkflow, err := s.repo.GetN8NWorkflowByID(workflowID)
	if err != nil {
		return "", fmt.Errorf("workflow not found: %w", err)
	}

	webhookID, _ := extractWebhookIDFromWorkflow(dbWorkflow.WorkflowJSON)
	if webhookID == "" {
		return "", fmt.Errorf("webhook ID not found in workflow")
	}

	webhookURL := fmt.Sprintf("%s/webhook/%s/webhook", N8N_BASE_URL, webhookID)
	if err := s.setTelegramWebhook(botToken, webhookURL); err != nil {
		return "", fmt.Errorf("setWebhook failed: %w", err)
	}

	log.Printf("[RegisterWebhook] workflow=%d webhookURL=%s", workflowID, webhookURL)
	return webhookURL, nil
}

func extractWebhookIDFromWorkflow(workflowJSON []byte) (string, error) {
	var wfData map[string]interface{}
	if err := json.Unmarshal(workflowJSON, &wfData); err != nil {
		return "", err
	}
	nodes, ok := wfData["nodes"].([]interface{})
	if !ok {
		return "", fmt.Errorf("no nodes")
	}
	for _, node := range nodes {
		nodeMap, ok := node.(map[string]interface{})
		if !ok {
			continue
		}
		if nodeMap["type"] != "n8n-nodes-base.telegramTrigger" {
			continue
		}
		if id, ok := nodeMap["webhookId"].(string); ok && id != "" {
			return id, nil
		}
	}
	return "", fmt.Errorf("webhookId not found")
}

func (s *Service) GetUserWorkflows(userID uint) ([]models.N8NWorkflow, error) {
	return s.repo.GetN8NWorkflowsByUserID(userID)
}

func (s *Service) GetWorkflowByID(workflowID uint) (*models.N8NWorkflow, error) {
	return s.repo.GetN8NWorkflowByID(workflowID)
}

func (s *Service) DeleteWorkflow(workflowID uint) error {
	workflow, err := s.repo.GetN8NWorkflowByID(workflowID)
	if err == nil && workflow.WorkflowID != "" {
		s.deleteN8NWorkflowByID(workflow.WorkflowID)
	}
	return s.repo.DeleteN8NWorkflow(workflowID)
}

func (s *Service) convertMapToN8NWorkflow(data map[string]interface{}, fallbackName string) *N8NWorkflowRequest {
	workflow := &N8NWorkflowRequest{
		Name:        fallbackName,
		Nodes:       []interface{}{},
		Connections: map[string]interface{}{},
		Settings:    map[string]interface{}{},
	}

	if name, ok := data["name"].(string); ok && name != "" {
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
	if len(workflow.Settings) == 0 {
		workflow.Settings = map[string]interface{}{"executionOrder": "v1"}
	}

	return workflow
}

func (s *Service) deleteN8NWorkflowByID(workflowID string) error {
	return nil
}

func injectWebhookIDs(data map[string]interface{}) string {
	webhookID := uuid.New().String()
	nodes, ok := data["nodes"].([]interface{})
	if !ok {
		return webhookID
	}
	for _, node := range nodes {
		nodeMap, ok := node.(map[string]interface{})
		if !ok {
			continue
		}
		if nodeMap["type"] != "n8n-nodes-base.telegramTrigger" {
			continue
		}
		nodeMap["webhookId"] = webhookID

		if params, ok := nodeMap["parameters"].(map[string]interface{}); ok {
			delete(params, "webhookId")
		}
	}
	return webhookID
}

func extractCredentialFromWorkflow(workflowJSON []byte) (credID, credName string) {
	credID = "bot-cred"
	credName = "TelegramBot"

	var wfData map[string]interface{}
	if err := json.Unmarshal(workflowJSON, &wfData); err != nil {
		return
	}

	nodes, ok := wfData["nodes"].([]interface{})
	if !ok {
		return
	}

	for _, node := range nodes {
		nodeMap, ok := node.(map[string]interface{})
		if !ok {
			continue
		}
		creds, ok := nodeMap["credentials"].(map[string]interface{})
		if !ok {
			continue
		}
		tgCred, ok := creds["telegramApi"].(map[string]interface{})
		if !ok {
			continue
		}
		if id, ok := tgCred["id"].(string); ok && id != "" {
			credID = id
		}
		if name, ok := tgCred["name"].(string); ok && name != "" {
			credName = name
		}
		return
	}
	return
}