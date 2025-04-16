package handlers

import (
	"fmt"
	"strings"
	"time"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/handlers/callbacks"
	"github.com/skyzeper/telegram-bot/internal/menus"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/security"
	"github.com/skyzeper/telegram-bot/internal/services/chat"
	"github.com/skyzeper/telegram-bot/internal/services/notification"
	"github.com/skyzeper/telegram-bot/internal/services/order"
	"github.com/skyzeper/telegram-bot/internal/services/user"
	"github.com/skyzeper/telegram-bot/internal/state"
)

// Handler manages incoming Telegram updates
type Handler struct {
	bot                *tgbotapi.BotAPI
	security           *security.SecurityChecker
	menus              *menus.MenuGenerator
	userService        *user.Service
	orderService       *order.Service
	chatService        *chat.Service
	state              *state.Manager
	callbackHandler    *callbacks.CallbackHandler
	notificationService *notification.Service
}

// NewHandler creates a new Handler
func NewHandler(
	bot *tgbotapi.BotAPI,
	security *security.SecurityChecker,
	menus *menus.MenuGenerator,
	userService *user.Service,
	orderService *order.Service,
	chatService *chat.Service,
	state *state.Manager,
	callbackHandler *callbacks.CallbackHandler,
	notificationService *notification.Service,
) *Handler {
	return &Handler{
		bot:                bot,
		security:           security,
		menus:              menus,
		userService:        userService,
		orderService:       orderService,
		chatService:        chatService,
		state:              state,
		callbackHandler:    callbackHandler,
		notificationService: notificationService,
	}
}

// HandleUpdate processes incoming Telegram updates
func (h *Handler) HandleUpdate(update *tgbotapi.Update) {
	if update.Message == nil {
		if update.CallbackQuery != nil {
			h.callbackHandler.HandleCallback(update.CallbackQuery)
		}
		return
	}

	chatID := update.Message.Chat.ID
	currentState := h.state.Get(chatID)

	// Check if user exists, create if not
	user, err := h.userService.GetUser(chatID)
	if err != nil {
		err = h.security.CreateUser(
			chatID,
			"client",
			update.Message.From.FirstName,
			update.Message.From.LastName,
			update.Message.From.UserName,
			"",
		)
		if err != nil {
			h.sendMessage(chatID, "❌ Ошибка регистрации. Попробуйте позже.", nil)
			return
		}
		user = &models.User{ChatID: chatID, Role: "client"}
	}

	// Handle commands
	if update.Message.IsCommand() {
		h.handleCommand(update, user)
		return
	}

	// Handle state-based interactions
	if currentState.Module != "" {
		switch currentState.Module {
		case "order":
			// Delegate to order steps handler
			// Note: This is handled in order/steps.go
			return
		case "chat":
			h.handleChatMessage(update)
			return
		}
	}

	// Handle text messages
	h.handleTextMessage(update, user)
}

// handleCommand processes Telegram commands
func (h *Handler) handleCommand(update *tgbotapi.Update, user *models.User) {
	chatID := update.Message.Chat.ID
	command := update.Message.Command()

	switch command {
	case "start":
		h.sendMessage(chatID, "Добро пожаловать! 🚛 Выберите действие:", h.menus.MainMenu(user))
	case "help":
		h.sendMessage(chatID, "📚 Помощь: Используйте меню для заказа услуг, связи с оператором или приглашения друзей.", h.menus.MainMenu(user))
	default:
		h.sendMessage(chatID, "❓ Неизвестная команда. Используйте /start или /help.", h.menus.MainMenu(user))
	}
}

// handleTextMessage processes text messages
func (h *Handler) handleTextMessage(update *tgbotapi.Update, user *models.User) {
	chatID := update.Message.Chat.ID
	messageText := strings.ToLower(update.Message.Text)

	switch messageText {
	case "🗑️ заказать услугу":
		role, err := h.security.GetUserRole(chatID)
		if err != nil {
			h.sendMessage(chatID, "❌ Ошибка проверки доступа. Попробуйте позже.", nil)
			return
		}
		if role == "client" || role == "operator" || role == "main_operator" || role == "owner" {
			h.state.Set(chatID, state.State{
				Module:     "order",
				Step:       1,
				TotalSteps: 11,
				Data:       make(map[string]interface{}),
			})
			h.sendMessage(chatID, "🗑️ Выберите категорию заказа:", h.menus.CategoryMenu())
		} else {
			h.sendMessage(chatID, "❌ У вас нет доступа к заказам.", nil)
		}

	case "📞 связаться с оператором":
		role, err := h.security.GetUserRole(chatID)
		if err != nil {
			h.sendMessage(chatID, "❌ Ошибка проверки доступа. Попробуйте позже.", nil)
			return
		}
		if role == "client" || role == "operator" || role == "main_operator" || role == "owner" {
			h.state.Set(chatID, state.State{
				Module:     "chat",
				Step:       1,
				TotalSteps: 1,
				Data:       make(map[string]interface{}),
			})
			h.sendMessage(chatID, "💬 Напишите ваш вопрос, и оператор ответит вам:", nil)
		} else {
			h.sendMessage(chatID, "❌ У вас нет доступа к чату.", nil)
		}

	case "🔗 приглашайте друзей и зарабатывайте!":
		role, err := h.security.GetUserRole(chatID)
		if err != nil {
			h.sendMessage(chatID, "❌ Ошибка проверки доступа. Попробуйте позже.", nil)
			return
		}
		if role == "client" || role == "operator" || role == "main_operator" || role == "owner" {
			h.sendMessage(chatID, "🔗 Поделитесь ссылкой для приглашения друзей:", h.menus.ReferralMenu())
		} else {
			h.sendMessage(chatID, "❌ У вас нет доступа к реферальной программе.", nil)
		}

	case "📋 заказы":
		role, err := h.security.GetUserRole(chatID)
		if err != nil {
			h.sendMessage(chatID, "❌ Ошибка проверки доступа. Попробуйте позже.", nil)
			return
		}
		if role == "operator" || role == "main_operator" || role == "owner" {
			h.sendMessage(chatID, "📋 Выберите категорию заказов:", h.menus.NewOrdersMenu(nil))
		} else {
			h.sendMessage(chatID, "❌ У вас нет доступа к заказам.", nil)
		}

	case "🧑‍💼 управление штатом":
		role, err := h.security.GetUserRole(chatID)
		if err != nil {
			h.sendMessage(chatID, "❌ Ошибка проверки доступа. Попробуйте позже.", nil)
			return
		}
		if role == "main_operator" || role == "owner" {
			h.sendMessage(chatID, "🧑‍💼 Управление штатом:", h.menus.StaffMenu(chatID))
		} else {
			h.sendMessage(chatID, "❌ У вас нет доступа к управлению штатом.", nil)
		}

	case "📊 статистика":
		role, err := h.security.GetUserRole(chatID)
		if err != nil {
			h.sendMessage(chatID, "❌ Ошибка проверки доступа. Попробуйте позже.", nil)
			return
		}
		if role == "owner" {
			h.sendMessage(chatID, "📊 Выберите период статистики:", nil)
		} else {
			h.sendMessage(chatID, "❌ У вас нет доступа к статистике.", nil)
		}

	default:
		h.sendMessage(chatID, "❓ Пожалуйста, выберите действие из меню:", h.menus.MainMenu(user))
	}
}

// handleChatMessage processes messages in chat mode
func (h *Handler) handleChatMessage(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	messageText := update.Message.Text

	// Save user message
	err := h.chatService.SaveMessage(chatID, messageText, true)
	if err != nil {
		h.sendMessage(chatID, "❌ Ошибка отправки сообщения. Попробуйте позже.", nil)
		return
	}

	// Notify operator
	operatorID, err := h.chatService.GetActiveOperator()
	if err != nil {
		h.sendMessage(chatID, "❌ Нет доступных операторов. Попробуйте позже.", nil)
		return
	}

	notification := &models.Notification{
		UserID:    operatorID,
		Type:      "chat_message",
		Message:   fmt.Sprintf("Новое сообщение от пользователя %d: %s", chatID, messageText),
		CreatedAt: time.Now(),
	}
	err = h.notificationService.CreateNotification(notification)
	if err != nil {
		h.sendMessage(chatID, "❌ Ошибка уведомления оператора. Попробуйте позже.", nil)
		return
	}

	h.sendMessage(chatID, "✅ Сообщение отправлено! Оператор скоро ответит.", nil)
}

// sendMessage sends a message to a chat
func (h *Handler) sendMessage(chatID int64, text string, replyMarkup interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	if replyMarkup != nil {
		switch rm := replyMarkup.(type) {
		case tgbotapi.ReplyKeyboardMarkup:
			msg.ReplyMarkup = rm
		case tgbotapi.InlineKeyboardMarkup:
			msg.ReplyMarkup = rm
		}
	}
	if _, err := h.bot.Send(msg); err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
	}
}