package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"main/models"
	"time"
)

// GenerateBotID generates a unique bot ID with "bot_" prefix
func GenerateBotID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "bot_" + hex.EncodeToString(bytes), nil
}

// CreateBot creates a new bot for the user
func (s *Service) CreateBot(userID, name, botType, productDescription string, faq []models.FAQItem) (*models.Bot, error) {
	// Validate input
	if name == "" {
		return nil, errors.New("bot name is required")
	}
	if botType == "" {
		return nil, errors.New("bot type is required")
	}

	// Generate unique bot ID
	botID, err := GenerateBotID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate bot ID: %w", err)
	}

	// Build system prompt from productDescription and FAQ
	systemPrompt := buildSystemPrompt(productDescription, faq)

	// Create bot
	bot := &models.Bot{
		BotID:        botID,
		UserID:       userID,
		Name:         name,
		Type:         botType,
		Status:       "active",
		SystemPrompt: systemPrompt,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateBot(bot); err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return bot, nil
}

// buildSystemPrompt constructs the AI system prompt from product description and FAQ
func buildSystemPrompt(productDescription string, faq []models.FAQItem) string {
	prompt := "Ты - AI ассистент компании, который помогает клиентам с вопросами о продуктах и услугах.\n\n"

	if productDescription != "" {
		prompt += "## Описание продуктов/услуг компании:\n"
		prompt += productDescription + "\n\n"
	}

	if len(faq) > 0 {
		prompt += "## Часто задаваемые вопросы (FAQ):\n\n"
		for i, item := range faq {
			prompt += fmt.Sprintf("%d. **Вопрос:** %s\n", i+1, item.Question)
			prompt += fmt.Sprintf("   **Ответ:** %s\n\n", item.Answer)
		}
	}

	prompt += "Используй эту информацию для ответов на вопросы клиентов. Будь вежливым, профессиональным и полезным."

	return prompt
}

// ProcessChatMessage processes a chat message and returns AI response
func (s *Service) ProcessChatMessage(botID, sessionID, message string) (string, error) {
	// Validate input
	if botID == "" {
		return "", errors.New("botId is required")
	}
	if sessionID == "" {
		return "", errors.New("sessionId is required")
	}
	if message == "" {
		return "", errors.New("message is required")
	}

	// Get bot to retrieve system prompt
	bot, err := s.repo.GetBotByBotID(botID)
	if err != nil {
		return "", errors.New("bot not found or inactive")
	}

	// Get chat history
	history, err := s.repo.GetSessionMessages(sessionID, botID)
	if err != nil {
		return "", fmt.Errorf("failed to get chat history: %w", err)
	}

	// Add user message to history
	userMessage := models.ChatMessage{
		Role:      "user",
		Content:   message,
		Timestamp: time.Now(),
	}
	if err := s.repo.AddMessageToSession(sessionID, botID, userMessage); err != nil {
		return "", fmt.Errorf("failed to save user message: %w", err)
	}

	// TODO: Integrate with AI (OpenAI/Anthropic)
	// For now, return a mock response
	aiReply := s.generateMockAIResponse(message, history, bot.SystemPrompt)

	// Save AI response to history
	aiMessage := models.ChatMessage{
		Role:      "assistant",
		Content:   aiReply,
		Timestamp: time.Now(),
	}
	if err := s.repo.AddMessageToSession(sessionID, botID, aiMessage); err != nil {
		return "", fmt.Errorf("failed to save AI message: %w", err)
	}

	return aiReply, nil
}

// generateMockAIResponse generates a mock AI response
// TODO: Replace with actual AI integration (OpenAI, Anthropic, etc.)
func (s *Service) generateMockAIResponse(message string, history []models.ChatMessage, systemPrompt string) string {
	// Simple mock response
	// When integrating real AI, pass systemPrompt as the system message
	return "Спасибо за сообщение! Это тестовый ответ от AI агента. Вы написали: \"" + message + "\""
}

// GetBotByID retrieves a bot by its ID
func (s *Service) GetBotByID(botID string) (*models.Bot, error) {
	return s.repo.GetBotByBotID(botID)
}

// GetUserBots retrieves all bots for a user
func (s *Service) GetUserBots(userID string) ([]models.Bot, error) {
	return s.repo.GetBotsByUserID(userID)
}

// ValidateBotOwnership checks if the bot belongs to the user
func (s *Service) ValidateBotOwnership(botID, userID string) error {
	return s.repo.ValidateBotOwnership(botID, userID)
}

// DeleteBot deletes a bot by its ID (database primary key)
func (s *Service) DeleteBot(agentID uint, userID string) error {
	return s.repo.DeleteBot(agentID, userID)
}
