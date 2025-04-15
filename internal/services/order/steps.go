package order

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) HandleOrderStepsWrapper(chatID int64, msg *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	s.HandleOrderSteps(chatID, msg, bot)
}
