package stats

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

func (s *Service) GetStats(ctx context.Context, period string, bot *tgbotapi.BotAPI, chatID int64) {
	// Placeholder for stats logic
	bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Статистика за %s", period)))
}
