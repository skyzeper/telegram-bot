package chat

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/utils"
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
	chatID := msg.Chat.ID

	if msg.Text == "Сам позвоню" {
		bot.Send(tgbotapi.NewMessage(chatID, "📞 Позвоните: +7(978)-959-70-77"))
		return
	}
	if msg.Text == "Позвоните мне" {
		msgConfig := tgbotapi.NewMessage(chatID, "Введите номер или нажмите 📞")
		msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📞 Отправить номер", "chat_contact"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "action_back"),
			),
		)
		bot.Send(msgConfig)
		return
	}
	if msg.Contact != nil || utils.FormatPhone(msg.Text) != "" {
		phone := msg.Contact.PhoneNumber
		if phone == "" {
			phone = utils.FormatPhone(msg.Text)
		}
		s.repo.SaveMessage(ctx, &models.Message{
			UserID:     chatID,
			Message:    fmt.Sprintf("Клиент %d просит позвонить: %s", chatID, phone),
			IsFromUser: true,
		})
		bot.Send(tgbotapi.NewMessage(chatID, "Оператор свяжется с вами"))
		return
	}
	s.repo.SaveMessage(ctx, &models.Message{
		UserID:     chatID,
		Message:    msg.Text,
		IsFromUser: true,
	})
	bot.Send(tgbotapi.NewMessage(chatID, "Сообщение отправлено оператору"))
}
