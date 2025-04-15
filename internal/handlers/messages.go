package handlers

import (
	"context"
	"database/sql"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/menus"
	"github.com/skyzeper/telegram-bot/internal/services/chat"
	"github.com/skyzeper/telegram-bot/internal/services/order"
	"github.com/skyzeper/telegram-bot/internal/services/referral"
	"github.com/skyzeper/telegram-bot/internal/services/user"
	"github.com/skyzeper/telegram-bot/internal/state"
)

type Handler struct {
	users    *user.Service
	orders   *order.Service
	chat     *chat.Service
	referral *referral.Service
	states   *state.StateManager
}

func NewHandler(db *sql.DB, states *state.StateManager) *Handler {
	return &Handler{
		users:    user.NewService(db, states),
		orders:   order.NewService(db, states),
		chat:     chat.NewService(db),
		referral: referral.NewService(db),
		states:   states,
	}
}

func HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message, db *sql.DB) {
	handler := NewHandler(db, state.NewStateManager())
	user, err := handler.users.GetUser(ctx, msg.Chat.ID)
	if err != nil && msg.Text != "/start" {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "🚫 Пожалуйста, начните с /start"))
		return
	}

	if msg.Text == "/start" {
		if err != nil {
			handler.users.Register(ctx, msg.Chat.ID, msg.From.FirstName, msg.From.LastName)
		}
		var menu tgbotapi.InlineKeyboardMarkup
		switch user.Role {
		case "user":
			menu = menus.UserMenu()
		case "operator":
			menu = menus.OperatorMenu()
		case "main_operator":
			menu = menus.MainOperatorMenu()
		case "driver":
			menu = menus.DriverMenu()
		case "loader":
			menu = menus.LoaderMenu()
		case "owner":
			menu = menus.OwnerMenu()
		}
		msgConfig := tgbotapi.NewMessage(msg.Chat.ID, "🚛 Заказ за 30 минут! 🔥 Лучшая цена! 😎 Простое оформление!")
		msgConfig.ReplyMarkup = menu
		bot.Send(msgConfig)
		return
	}

	currentState := handler.states.Get(msg.Chat.ID)
	if currentState.Module != "" {
		switch currentState.Module {
		case "create_order":
			handler.orders.HandleOrderSteps(ctx, msg.Chat.ID, msg, bot)
		case "add_staff":
			handler.users.HandleStaffSteps(msg.Chat.ID, msg, bot)
		case "chat":
			handler.chat.HandleChat(ctx, msg, bot)
		}
		return
	}

	// Если сообщение не команда и не часть шага, перенаправляем в чат
	handler.chat.HandleChat(ctx, msg, bot)
}
