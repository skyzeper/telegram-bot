package order

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) HandleOrderStepsWrapper(ctx context.Context, chatID int64, msg *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	s.HandleOrderSteps(ctx, chatID, msg, bot)
}
