package payment

import (
	"context"
	"database/sql"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(db *sql.DB) *Service {
	return &Service{
		repo: NewRepository(db),
	}
}

func (s *Service) CreatePayment(ctx context.Context, orderID, userID int64, amount float64, method string) error {
	payment := &models.Payment{
		OrderID: orderID,
		UserID:  userID,
		Amount:  amount,
		Method:  method,
	}
	return s.repo.CreatePayment(ctx, payment)
}

func (s *Service) ConfirmPayment(ctx context.Context, paymentID int64, bot *tgbotapi.BotAPI) error {
	payment, err := s.repo.GetPayment(ctx, paymentID)
	if err != nil {
		return err
	}
	payment.Confirmed = true
	if err := s.repo.UpdatePayment(ctx, payment); err != nil {
		return err
	}
	bot.Send(tgbotapi.NewMessage(payment.UserID, "Оплата подтверждена"))
	return nil
}
