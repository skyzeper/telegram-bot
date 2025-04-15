package user

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) HandleStaffSteps(chatID int64, msg *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	currentState := s.states.Get(chatID)
	if currentState.Module != "add_staff" {
		return
	}
	s.AddStaff(context.Background(), chatID, bot, msg)
}
