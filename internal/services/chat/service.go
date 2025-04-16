package chat

import (
	"errors"
	"fmt"
	"time"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/models"
)

// Service handles chat-related business logic
type Service struct {
	bot  *tgbotapi.BotAPI
	repo Repository
}

// Repository defines the interface for chat data access
type Repository interface {
	CreateMessage(message *models.Message) error
	GetMessagesByUser(userID int64) ([]models.Message, error)
	GetActiveOperator() (int64, error)
}

// NewService creates a new chat service
func NewService(bot *tgbotapi.BotAPI, repo Repository) *Service {
	return &Service{
		bot:  bot,
		repo: repo,
	}
}

// StartChat initiates a chat session
func (s *Service) StartChat(userID int64) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}

	// Notify operators
	operatorID, err := s.repo.GetActiveOperator()
	if err != nil {
		return fmt.Errorf("failed to get active operator: %v", err)
	}

	if operatorID == 0 {
		return errors.New("no active operators available")
	}

	notifyMsg := tgbotapi.NewMessage(operatorID, fmt.Sprintf(
		"💬 Новый чат начат с пользователем (Chat ID: %d).",
		userID,
	))
	if _, err := s.bot.Send(notifyMsg); err != nil {
		return fmt.Errorf("failed to notify operator: %v", err)
	}

	return nil
}

// ForwardMessageToOperator forwards a user message to the operator
func (s *Service) ForwardMessageToOperator(msg *tgbotapi.Message) error {
	if msg.Chat.ID <= 0 || msg.Text == "" {
		return errors.New("invalid message or chat ID")
	}

	operatorID, err := s.repo.GetActiveOperator()
	if err != nil {
		return fmt.Errorf("failed to get active operator: %v", err)
	}

	if operatorID == 0 {
		return errors.New("no active operators available")
	}

	forwardMsg := tgbotapi.NewMessage(operatorID, fmt.Sprintf(
		"> Сообщение от пользователя (Chat ID: %d):\n%s",
		msg.Chat.ID,
		msg.Text,
	))
	forwardMsg.ParseMode = "Markdown"
	if _, err := s.bot.Send(forwardMsg); err != nil {
		return fmt.Errorf("failed to forward message: %v", err)
	}

	// Save message
	message := &models.Message{
		UserID:     msg.Chat.ID,
		OperatorID: operatorID,
		Message:    msg.Text,
		IsFromUser: true,
		CreatedAt:  time.Now(),
	}
	return s.repo.CreateMessage(message)
}

// IsInChat checks if a user is in a chat session
func (s *Service) IsInChat(userID int64) bool {
	// Placeholder: In a real implementation, check active chat sessions
	// For mock, assume user is in chat if they sent a message recently
	messages, err := s.repo.GetMessagesByUser(userID)
	if err != nil {
		return false
	}
	for _, msg := range messages {
		if time.Since(msg.CreatedAt) < 24*time.Hour {
			return true
		}
	}
	return false
}