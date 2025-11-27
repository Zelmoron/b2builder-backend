package repository

import (
	"encoding/json"
	"errors"
	"main/models"
	"time"
)

// CreateBot creates a new bot in the database
func (r *Repository) CreateBot(bot *models.Bot) error {
	result := r.db.Create(bot)
	return result.Error
}

// GetBotByBotID retrieves a bot by its bot_id
func (r *Repository) GetBotByBotID(botID string) (*models.Bot, error) {
	var bot models.Bot
	result := r.db.Where("bot_id = ? AND status = ?", botID, "active").First(&bot)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bot, nil
}

// GetBotsByUserID retrieves all bots for a specific user
func (r *Repository) GetBotsByUserID(userID string) ([]models.Bot, error) {
	var bots []models.Bot
	result := r.db.Where("user_id = ?", userID).Find(&bots)
	if result.Error != nil {
		return nil, result.Error
	}
	return bots, nil
}

// GetChatSession retrieves a chat session by session_id and bot_id
func (r *Repository) GetChatSession(sessionID, botID string) (*models.ChatSession, error) {
	var session models.ChatSession
	result := r.db.Where("session_id = ? AND bot_id = ?", sessionID, botID).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

// CreateChatSession creates a new chat session
func (r *Repository) CreateChatSession(session *models.ChatSession) error {
	result := r.db.Create(session)
	return result.Error
}

// UpdateChatSession updates an existing chat session
func (r *Repository) UpdateChatSession(session *models.ChatSession) error {
	result := r.db.Save(session)
	return result.Error
}

// AddMessageToSession adds a message to a chat session
func (r *Repository) AddMessageToSession(sessionID, botID string, message models.ChatMessage) error {
	session, err := r.GetChatSession(sessionID, botID)
	if err != nil {
		// Create new session if it doesn't exist
		messages := []models.ChatMessage{message}
		messagesJSON, err := json.Marshal(messages)
		if err != nil {
			return err
		}

		newSession := &models.ChatSession{
			SessionID: sessionID,
			BotID:     botID,
			Messages:  messagesJSON,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return r.CreateChatSession(newSession)
	}

	// Unmarshal existing messages
	var messages []models.ChatMessage
	if err := json.Unmarshal(session.Messages, &messages); err != nil {
		return err
	}

	// Add new message
	messages = append(messages, message)

	// Marshal back to JSON
	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	session.Messages = messagesJSON
	session.UpdatedAt = time.Now()

	return r.UpdateChatSession(session)
}

// GetSessionMessages retrieves all messages from a session
func (r *Repository) GetSessionMessages(sessionID, botID string) ([]models.ChatMessage, error) {
	session, err := r.GetChatSession(sessionID, botID)
	if err != nil {
		return []models.ChatMessage{}, nil // Return empty array if session doesn't exist
	}

	var messages []models.ChatMessage
	if err := json.Unmarshal(session.Messages, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// BotExists checks if a bot with the given bot_id exists and is active
func (r *Repository) BotExists(botID string) (bool, error) {
	var count int64
	result := r.db.Model(&models.Bot{}).Where("bot_id = ? AND status = ?", botID, "active").Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// ValidateBotOwnership checks if the bot belongs to the user
func (r *Repository) ValidateBotOwnership(botID, userID string) error {
	var bot models.Bot
	result := r.db.Where("bot_id = ? AND user_id = ?", botID, userID).First(&bot)
	if result.Error != nil {
		return errors.New("bot not found or access denied")
	}
	return nil
}