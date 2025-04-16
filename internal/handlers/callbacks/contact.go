package callbacks

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/menus"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/security"
	"github.com/skyzeper/telegram-bot/internal/services/chat"
	"github.com/skyzeper/telegram-bot/internal/state"
	"github.com/skyzeper/telegram-bot/internal/utils"
)

// ContactHandler handles contact-related callback queries
type ContactHandler struct {
	bot           *tgbotapi.BotAPI
	security      *security.SecurityChecker
	menus         *menus.MenuGenerator
	chatService   *chat.Service
	state         *state.Manager
}

// NewContactHandler creates a new ContactHandler
func NewContactHandler(
	bot *tgbotapi.BotAPI,
	security *security.SecurityChecker,
	menus *menus.MenuGenerator,
	chatService *chat.Service,
	state *state.Manager,
) *ContactHandler {
	return &ContactHandler{
		bot:           bot,
		security:      security,
		menus:         menus,
		chatService:   chatService,
		state:         state,
	}
}

// Handle processes contact-related callbacks
func (h *ContactHandler) Handle(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	switch data {
	case "contact_call":
		h.sendMessage(chatID, callback.Message.MessageID, "📞 Пожалуйста, свяжитесь с нами по телефону: +1234567890")
	case "contact_request_call":
		h.sendMessage(chatID, callback.Message.MessageID, "📲 Мы свяжемся с вами в ближайшее время!")
		// Notify operator
		operatorID, err := h.chatService.GetActiveOperator()
		if err == nil {
			h.chatService.SaveMessage(operatorID, fmt.Sprintf("Пользователь %d запрашивает звонок", chatID), false)
		}
	case "contact_chat":
		user := &models.User{ChatID: chatID}
		h.state.Set(chatID, state.State{
			Module:     "chat",
			Step:       1,
			TotalSteps: 1,
			Data:       make(map[string]interface{}),
		})
		h.sendMessage(chatID, callback.Message.MessageID, "💬 Напишите ваш вопрос, и оператор ответит вам:")
	default:
		h.sendMessage(chatID, callback.Message.MessageID, "❓ Неизвестная команда.")
	}
}

// sendMessage sends a message in response to a callback
func (h *ContactHandler) sendMessage(chatID int64, messageID int, text string, replyMarkup ...interface{}) {
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