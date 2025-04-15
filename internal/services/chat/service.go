package chat

import (
	"context"
	"database/sql"
	"fmt"

	"bot/internal/models"
	"bot/internal/utils"
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

func (s *Service) HandleChat(ctx context.Context, msg *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if msg.Text == "Сам позвоню" {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "📞 Позвоните: +7(978)-959-70-77"))
		return
	}
	if msg.Text == "Позвоните мне" {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Введите номер или нажмите 📞", tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonContact("📞 Отправить номер"),
			),
		)))
		return
	}
	if msg.Contact != nil || utils.FormatPhone(msg.Text) != "" {
		phone := msg.Contact.PhoneNumber
		if phone == "" {
			phone = utils.FormatPhone(msg.Text)
		}
		s.repo.SaveMessage(ctx, &models.Message{
			UserID:     msg.Chat.ID,
			Message:    fmt.Sprintf("Клиент %d просит позвонить: %s", msg.Chat.ID, phone),
			IsFromUser: true,
		})
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Оператор свяжется с вами", tgbotapi.ReplyKeyboardRemove{}))
		return
	}
	s.repo.SaveMessage(ctx, &models.Message{
		UserID:     msg.Chat.ID,
		Message:    msg.Text,
		IsFromUser: true,
	})
	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Сообщение отправлено оператору"))
}
