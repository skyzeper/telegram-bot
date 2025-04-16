package callbacks

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/menus"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/security"
	"github.com/skyzeper/telegram-bot/internal/services/chat"
	"github.com/skyzeper/telegram-bot/internal/services/executor"
	"github.com/skyzeper/telegram-bot/internal/services/order"
	"github.com/skyzeper/telegram-bot/internal/services/payment"
	"github.com/skyzeper/telegram-bot/internal/services/review"
	"github.com/skyzeper/telegram-bot/internal/services/user"
	"github.com/skyzeper/telegram-bot/internal/state"
	"github.com/skyzeper/telegram-bot/internal/utils"
)

// OrdersHandler handles order-related callback queries
type OrdersHandler struct {
	bot           *tgbotapi.BotAPI
	security      *security.SecurityChecker
	menus         *menus.MenuGenerator
	userService   *user.Service
	orderService  *order.Service
	chatService   *chat.Service
	executorService *executor.Service
	paymentService *payment.Service
	reviewService *review.Service
	state         *state.Manager
}

// NewOrdersHandler creates a new OrdersHandler
func NewOrdersHandler(
	bot *tgbotapi.BotAPI,
	security *security.SecurityChecker,
	menus *menus.MenuGenerator,
	userService *user.Service,
	orderService *order.Service,
	chatService *chat.Service,
	executorService *executor.Service,
	paymentService *payment.Service,
	reviewService *review.Service,
	state *state.Manager,
) *OrdersHandler {
	return &OrdersHandler{
		bot:           bot,
		security:      security,
		menus:         menus,
		userService:   userService,
		orderService:  orderService,
		chatService:   chatService,
		executorService: executorService,
		paymentService: paymentService,
		reviewService: reviewService,
		state:         state,
	}
}

// Handle processes order-related callbacks
func (h *OrdersHandler) Handle(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	// Placeholder for order handling logic
	switch data {
	case "accept_order_1":
		h.sendMessage(chatID, callback.Message.MessageID, "✅ Заказ принят!")
	case "cancel_order_1":
		h.sendMessage(chatID, callback.Message.MessageID, "❌ Заказ отменён.")
	default:
		h.sendMessage(chatID, callback.Message.MessageID, "❓ Неизвестная команда.")
	}
}

// sendMessage sends a message in response to a callback
func (h *OrdersHandler) sendMessage(chatID int64, messageID int, text string, replyMarkup ...interface{}) {
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