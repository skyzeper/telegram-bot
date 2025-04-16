package callbacks

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/menus"
	"github.com/skyzeper/telegram-bot/internal/services/review"
	"github.com/skyzeper/telegram-bot/internal/state"
	"github.com/skyzeper/telegram-bot/internal/security"
	"github.com/skyzeper/telegram-bot/internal/utils"
)

// ReviewsHandler handles review-related callbacks
type ReviewsHandler struct {
	bot           *tgbotapi.BotAPI
	security      security.SecurityChecker
	menus         *menus.MenuGenerator
	reviewService *review.Service
	state         *state.State
}

// NewReviewsHandler creates a new ReviewsHandler
func NewReviewsHandler(
	bot *tgbotapi.BotAPI,
	security security.SecurityChecker,
	menus *menus.MenuGenerator,
	reviewService *review.Service,
	state *state.State,
) *ReviewsHandler {
	return &ReviewsHandler{
		bot:           bot,
		security:      security,
		menus:         menus,
		reviewService: reviewService,
		state:         state,
	}
}

// HandleReviewsCallback processes review-related callback queries
func (h *ReviewsHandler) HandleReviewsCallback(callback *tgbotapi.CallbackQuery, user *models.User, data string) {
	if strings.HasPrefix(data, "review_rate_") {
		h.handleReviewRate(callback, user, data)
	} else if strings.HasPrefix(data, "rate_") {
		h.handleRatingSubmit(callback, user, data)
	}
}

// handleReviewRate processes review rating
func (h *ReviewsHandler) handleReviewRate(callback *tgbotapi.CallbackQuery, user *models.User, data string) {
	if !h.security.HasRole(callback.Message.Chat.ID, "user") {
		h.sendUnauthorized(callback.Message.Chat.ID)
		return
	}

	orderIDStr := strings.TrimPrefix(data, "review_rate_")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		h.sendError(callback.Message.Chat.ID, "Неверный формат заказа.")
		return
	}

	h.state.SetState(callback.Message.Chat.ID, state.State{
		Module:      "review",
		Step:        1,
		TotalSteps:  2,
		Data:        map[string]interface{}{"order_id": orderID},
	})

	reply := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "🌟 Оцените заказ (1-5):")
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1", fmt.Sprintf("rate_%d_1", orderID)),
			tgbotapi.NewInlineKeyboardButtonData("2", fmt.Sprintf("rate_%d_2", orderID)),
			tgbotapi.NewInlineKeyboardButtonData("3", fmt.Sprintf("rate_%d_3", orderID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("4", fmt.Sprintf("rate_%d_4", orderID)),
			tgbotapi.NewInlineKeyboardButtonData("5", fmt.Sprintf("rate_%d_5", orderID)),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// handleRatingSubmit submits the review rating
func (h *ReviewsHandler) handleRatingSubmit(callback *tgbotapi.CallbackQuery, user *models.User, data string) {
	if !h.security.HasRole(callback.Message.Chat.ID, "user") {
		h.sendUnauthorized(callback.Message.Chat.ID)
		return
	}

	parts := strings.Split(strings.TrimPrefix(data, "rate_"), "_")
	if len(parts) != 2 {
		h.sendError(callback.Message.Chat.ID, "Неверный формат оценки.")
		return
	}

	orderID, err := strconv.Atoi(parts[0])
	if err != nil {
		h.sendError(callback.Message.Chat.ID, "Неверный формат заказа.")
		return
	}

	rating, err := strconv.Atoi(parts[1])
	if err != nil || rating < 1 || rating > 5 {
		h.sendError(callback.Message.Chat.ID, "Неверная оценка.")
		return
	}

	if err := h.reviewService.SubmitReview(orderID, user.ChatID, rating, ""); err != nil {
		h.sendError(callback.Message.Chat.ID, "Ошибка отправки отзыва.")
		return
	}

	h.state.ClearState(callback.Message.Chat.ID)

	reply := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		"🌟 Спасибо за ваш отзыв! 🙌\nВаш голос помогает нам стать лучше!",
	)
	reply.ReplyMarkup = h.menus.MainMenu(user)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// sendMainMenu sends the main menu
func (h *ReviewsHandler) sendMainMenu(chatID int64, user *models.User) {
	reply := tgbotapi.NewMessage(chatID, "Выберите действие:")
	reply.ReplyMarkup = h.menus.MainMenu(user)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// sendError sends an error message
func (h *ReviewsHandler) sendError(chatID int64, text string) {
	reply := tgbotapi.NewMessage(chatID, text)
	reply.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔙 Главное меню"),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// sendUnauthorized sends unauthorized access message
func (h *ReviewsHandler) sendUnauthorized(chatID int64) {
	reply := tgbotapi.NewMessage(chatID, "🚫 Доступ запрещён.")
	reply.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔙 Главное меню"),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}