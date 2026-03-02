package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	N8N_BASE_URL          = "https://prod.xoren8n.ru"
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

type N8NCredentialRequest struct {
	Name string                 `json:"name"`
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type N8NCredentialResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (s *Service) createN8NCredential(req *N8NCredentialRequest) (*N8NCredentialResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal credential: %w", err)
	}

	httpReq, err := http.NewRequest("POST", N8N_API_URL+"/credentials", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-N8N-API-KEY", N8N_API_KEY)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("[N8N Credential] status=%d body=%s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("n8n API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	credResp := &N8NCredentialResponse{
		Name: fmt.Sprintf("%v", raw["name"]),
		Type: fmt.Sprintf("%v", raw["type"]),
		ID:   fmt.Sprintf("%v", raw["id"]),
	}

	if credResp.ID == "" || credResp.ID == "<nil>" {
		return nil, fmt.Errorf("credential created but ID is empty, raw: %s", string(body))
	}

	return credResp, nil
}

func (s *Service) setTelegramWebhook(botToken, webhookURL string) error {
	type payload struct {
		URL                string `json:"url"`
		DropPendingUpdates bool   `json:"drop_pending_updates"`
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", botToken)
	client := &http.Client{}

	for attempt := 0; attempt < 3; attempt++ {
		body, _ := json.Marshal(payload{URL: webhookURL, DropPendingUpdates: true})
		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		log.Printf("[Telegram setWebhook] attempt=%d url=%s response=%s", attempt+1, webhookURL, string(respBody))

		if resp.StatusCode == http.StatusTooManyRequests {
			var result struct {
				Parameters struct {
					RetryAfter int `json:"retry_after"`
				} `json:"parameters"`
			}
			json.Unmarshal(respBody, &result)
			wait := result.Parameters.RetryAfter
			if wait <= 0 {
				wait = 2
			}
			log.Printf("[Telegram setWebhook] rate limited, retrying after %ds", wait)
			time.Sleep(time.Duration(wait) * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("telegram API error: status %d, body: %s", resp.StatusCode, string(respBody))
		}

		return nil
	}

	return fmt.Errorf("telegram setWebhook failed after retries")
}

func (s *Service) activateN8NWorkflow(workflowID string) error {
	url := fmt.Sprintf("%s/workflows/%s/activate", N8N_API_URL, workflowID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-N8N-API-KEY", N8N_API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("n8n API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}