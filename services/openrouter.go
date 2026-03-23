package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	OPENROUTER_API_URL = "https://openrouter.ai/api/v1/chat/completions"
	OPENROUTER_API_KEY = "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"
	GEMINI_MODEL       = "google/gemini-3-flash-preview"
)

type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterRequest struct {
	Model     string              `json:"model"`
	Messages  []OpenRouterMessage `json:"messages"`
	MaxTokens int                 `json:"max_tokens"`
	Reasoning map[string]bool     `json:"reasoning,omitempty"`
}

type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

func (s *Service) callOpenRouter(messages []OpenRouterMessage) (string, error) {
	reqBody := OpenRouterRequest{
		Model:     GEMINI_MODEL,
		Messages:  messages,
		MaxTokens: 8192,
		Reasoning: map[string]bool{
			"enabled": true,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", OPENROUTER_API_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", OPENROUTER_API_KEY))
	req.Header.Set("HTTP-Referer", "http://localhost:8080")
	req.Header.Set("X-Title", "B2B Builder")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenRouter API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if openRouterResp.Error != nil {
		return "", fmt.Errorf("OpenRouter error: %s", openRouterResp.Error.Message)
	}

	if len(openRouterResp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return openRouterResp.Choices[0].Message.Content, nil
}
