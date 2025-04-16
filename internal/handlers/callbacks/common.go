package callbacks

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/menus"
	"github.com/skyzeper/telegram-bot/internal/security"
	"github.com/skyzeper/telegram-bot/internal/services/user"
	"github.com/skyzeper/telegram-bot/internal/state"
	"github.com/skyzeper/telegram-bot/internal/utils"
	"strings"
)

// CallbackHandlable defines the interface for callback handlers
type CallbackHandlable interface {
	Handle(callback *tgbotapi.CallbackQuery)
}

// CallbackHandler manages callback queries
type CallbackHandler struct {
	bot              *tgbotapi.BotAPI
	security         *security.SecurityChecker
	menus            *menus.MenuGenerator
	userService      *user.Service
	state            *state.Manager
	ordersHandler    CallbackHandlable
	staffHandler     CallbackHandlable
	contactHandler   CallbackHandlable
	referralsHandler CallbackHandlable
	reviewsHandler   CallbackHandlable
	statsHandler     CallbackHandlable
}

// NewCallbackHandler creates a new CallbackHandler
func NewCallbackHandler(
	bot *tgbotapi.BotAPI,
	security *security.SecurityChecker,
	menus *menus.MenuGenerator,
	userService *user.Service,
	state *state.Manager,
	ordersHandler CallbackHandlable,
	staffHandler CallbackHandlable,
	contactHandler CallbackHandlable,
	referralsHandler CallbackHandlable,
	reviewsHandler CallbackHandlable,
	statsHandler CallbackHandlable,
) *CallbackHandler {
	return &CallbackHandler{
		bot:              bot,
		security:         security,
		menus:            menus,
		userService:      userService,
		state:            state,
		ordersHandler:    ordersHandler,
		staffHandler:     staffHandler,
		contactHandler:   contactHandler,
		referralsHandler: referralsHandler,
		reviewsHandler:   reviewsHandler,
		statsHandler:     statsHandler,
	}
}

// HandleCallback processes callback queries
func (h *CallbackHandler) HandleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	// Acknowledge callback
	h.bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	// Parse callback data
	parts := strings.Split(data, "_")
	if len(parts) < 1 {
		h.sendCallbackMessage(chatID, callback.Message.MessageID, "❌ Неверный формат команды.")
		return
	}

	module := parts[0]
	switch module {
	case "order", "accept", "cancel", "block", "assign", "confirm", "cash":
		h.ordersHandler.Handle(callback)
	case "staff", "edit":
		h.staffHandler.Handle(callback)
	case "contact", "call", "chat":
		h.contactHandler.Handle(callback)
	case "referral", "qr":
		h.referralsHandler.Handle(callback)
	case "review":
		h.reviewsHandler.Handle(callback)
	case "stats":
		h.statsHandler.Handle(callback)
	case "date":
		if len(parts) < 2 {
			h.sendCallbackMessage(chatID, callback.Message.MessageID, "❌ Неверный формат даты.")
			return
		}
		dateStr := parts[1]
		currentState := h.state.Get(chatID)
		if currentState.Module == "order" && currentState.Step == 3 {
			currentState.Data["date"] = dateStr
			currentState.Step = 4
			h.state.Set(chatID, currentState)
			h.sendCallbackMessage(chatID, callback.Message.MessageID, "🕒 Выберите время заказа:", h.menus.TimeMenu())
		}
	case "prev", "next":
		h.sendCallbackMessage(chatID, callback.Message.MessageID, "📅 Выберите дату:", h.menus.DateMenu())
	default:
		h.sendCallbackMessage(chatID, callback.Message.MessageID, "❓ Неизвестная команда.")
	}
}

// sendCallbackMessage sends a message in response to a callback
func (h *CallbackHandler) sendCallbackMessage(chatID int64, messageID int, text string, replyMarkup ...interface{}) {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	if len(replyMarkup) > 0 {
		if rm, ok := replyMarkup[0].(tgbotapi.InlineKeyboardMarkup); ok {
			msg.ReplyMarkup = &rm
		}
	}
	if _, err := h.bot.Send(msg); err != nil {
		utils.LogError(err)
	}
}
