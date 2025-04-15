package executor

import (
	"context"
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/services/notification"
)

type Service struct {
	repo   *Repository
	notify *notification.Service
}

func NewService(db *sql.DB, notify *notification.Service) *Service {
	return &Service{
		repo:   NewRepository(db),
		notify: notify,
	}
}

func (s *Service) AssignExecutor(ctx context.Context, orderID, userID int64, role string) error {
	executor := &models.Executor{
		OrderID: orderID,
		UserID:  userID,
		Role:    role,
	}
	if err := s.repo.CreateExecutor(ctx, executor); err != nil {
		return err
	}
	return s.notify.SendExecutorNotification(ctx, orderID, userID)
}

func (s *Service) ConfirmExecutor(ctx context.Context, orderID, userID int64, bot *tgbotapi.BotAPI) error {
	executor, err := s.repo.GetExecutor(ctx, orderID, userID)
	if err != nil {
		return err
	}
	executor.Confirmed = true
	if err := s.repo.UpdateExecutor(ctx, executor); err != nil {
		return err
	}
	bot.Send(tgbotapi.NewMessage(userID, "Заказ подтверждён"))
	return nil
}
