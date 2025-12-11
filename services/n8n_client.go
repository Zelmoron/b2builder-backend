package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	N8N_API_URL           = "https://prod.xoren8n.ru/api/v1"
	N8N_API_KEY           = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJiZTgxMmNhMC0zYWYyLTQ3ZWEtOTFhYy1iYTgzN2FlZGQzMmIiLCJpc3MiOiJuOG4iLCJhdWQiOiJwdWJsaWMtYXBpIiwiaWF0IjoxNzY1NDE2NDMxfQ.RjZ-CdDJRlWL1HvqKG3-xumbxZYrQ9XK2kvNd3zsmeo"
	DEFAULT_WORKFLOW_NAME = "Auto Generated Workflow"
)

type N8NWorkflowRequest struct {
	Name        string                 `json:"name"`
	Nodes       []interface{}          `json:"nodes"`
	Connections map[string]interface{} `json:"connections"`
	Active      bool                   `json:"active,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
}

type N8NWorkflowResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Active      bool                   `json:"active"`
	CreatedAt   string                 `json:"createdAt"`
	UpdatedAt   string                 `json:"updatedAt"`
	Nodes       []interface{}          `json:"nodes"`
	Connections map[string]interface{} `json:"connections"`
	Tags        []string               `json:"tags,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
}

func (s *Service) createN8NWorkflow(workflow *N8NWorkflowRequest) (*N8NWorkflowResponse, error) {
	jsonData, err := json.Marshal(workflow)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow: %w", err)
	}

	req, err := http.NewRequest("POST", N8N_API_URL+"/workflows", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-N8N-API-KEY", N8N_API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("n8n API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var workflowResp N8NWorkflowResponse
	if err := json.Unmarshal(body, &workflowResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &workflowResp, nil
}

func (s *Service) updateN8NWorkflowByID(workflowID string, workflow *N8NWorkflowRequest) (*N8NWorkflowResponse, error) {
	if workflowID == "" {
		return nil, fmt.Errorf("workflow ID is required")
	}

	jsonData, err := json.Marshal(workflow)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow: %w", err)
	}

	url := fmt.Sprintf("%s/workflows/%s", N8N_API_URL, workflowID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-N8N-API-KEY", N8N_API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("n8n API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var workflowResp N8NWorkflowResponse
	if err := json.Unmarshal(body, &workflowResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &workflowResp, nil
}

func (s *Service) getN8NWorkflowByID(workflowID string) (*N8NWorkflowResponse, error) {
	url := fmt.Sprintf("%s/workflows/%s", N8N_API_URL, workflowID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-N8N-API-KEY", N8N_API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("n8n API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var workflow N8NWorkflowResponse
	if err := json.Unmarshal(body, &workflow); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &workflow, nil
}