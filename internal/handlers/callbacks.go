package handlers

import (
	"context"
	"database/sql"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) HandleCallback(ctx context.Context, bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, db *sql.DB) {
	bot.Send(tgbotapi.NewMessage(callback.Message.Chat.ID, "Callback received: "+callback.Data))
}
