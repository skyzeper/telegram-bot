package handlers

import (
	"context"
	"database/sql"

	"bot/internal/menus"
	"bot/internal/services/user"
	"bot/internal/state"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	users  *user.Service
	states *state.StateManager
}

func NewHandler(db *sql.DB, states *state.StateManager) *Handler {
	return &Handler{
		users:  user.NewService(db, states),
		states: states,
	}
}

func HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, db *sql.DB) {
	handler := NewHandler(db, state.NewStateManager())
	switch msg.Text {
	case "/start":
		if _, err := handler.users.GetUser(ctx, msg.Chat.ID); err != nil {
			handler.users.Register(ctx, msg.Chat.ID, msg.From.FirstName, msg.From.LastName)
		}
		message := tgbotapi.NewMessage(msg.Chat.ID, "🚛 Заказ за 30 минут! 🔥 Лучшая цена! 😎 Простое оформление!")
		message.ReplyMarkup = menus.UserMenu()
		bot.Send(message)
	default:
		currentState := handler.states.Get(msg.Chat.ID)
		if currentState.Module == "add_staff" {
			handler.users.HandleStaffSteps(msg.Chat.ID, msg, bot)
		} else {
			message := tgbotapi.NewMessage(msg.Chat.ID, "Выберите действие")
			message.ReplyMarkup = menus.UserMenu()
			bot.Send(message)
		}
	}
}
