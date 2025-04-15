package notification

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	repo *Repository
}

func NewService(db *sql.DB) *Service {
	return &Service{
		repo: NewRepository(db),
	}
}

func (s *Service) SendExecutorNotification(ctx context.Context, orderID, userID int64) error {
	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Назначен заказ #%d. Подробности: ...", orderID))
	// Placeholder for actual order details
	return nil
}
