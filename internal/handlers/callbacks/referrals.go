package callbacks

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/menus"
	"github.com/skyzeper/telegram-bot/internal/services/referral"
	"github.com/skyzeper/telegram-bot/internal/services/user"
	"github.com/skyzeper/telegram-bot/internal/security"
	"github.com/skyzeper/telegram-bot/internal/utils"
)

// ReferralsHandler handles referral-related callbacks
type ReferralsHandler struct {
	bot             *tgbotapi.BotAPI
	security        security.SecurityChecker
	menus           *menus.MenuGenerator
	referralService *referral.Service
	userService     *user.Service
}

// NewReferralsHandler creates a new ReferralsHandler
func NewReferralsHandler(
	bot *tgbotapi.BotAPI,
	security security.SecurityChecker,
	menus *menus.MenuGenerator,
	referralService *referral.Service,
	userService *user.Service,
) *ReferralsHandler {
	return &ReferralsHandler{
		bot:             bot,
		security:        security,
		menus:           menus,
		referralService: referralService,
		userService:     userService,
	}
}

// HandleReferralsCallback processes referral-related callback queries
func (h *ReferralsHandler) HandleReferralsCallback(callback *tgbotapi.CallbackQuery, user *models.User, data string) {
	switch {
	case data == "referral_link":
		h.handleReferralLink(callback, user)
	case data == "referral_qr":
		h.handleReferralQR(callback, user)
	case strings.HasPrefix(data, "referral_payout_"):
		h.handleReferralPayout(callback, user, data)
	}
}

// handleReferralLink generates referral link
func (h *ReferralsHandler) handleReferralLink(callback *tgbotapi.CallbackQuery, user *models.User) {
	if !h.security.HasRole(callback.Message.Chat.ID, "user") {
		h.sendUnauthorized(callback.Message.Chat.ID)
		return
	}

	link := fmt.Sprintf("@vseVsimferopole?start=ref_%d", callback.Message.Chat.ID)
	reply := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		fmt.Sprintf("🔗 Ваша реферальная ссылка:\n%s\n\nПриглашайте друзей и получайте 500 рублей за заказ от 10,000 рублей! 🎉", link),
	)
	reply.ReplyMarkup = h.menus.ReferralMenu()
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// handleReferralQR generates QR code
func (h *ReferralsHandler) handleReferralQR(callback *tgbotapi.CallbackQuery, user *models.User) {
	if !h.security.HasRole(callback.Message.Chat.ID, "user") {
		h.sendUnauthorized(callback.Message.Chat.ID)
		return
	}

	link := fmt.Sprintf("@vseVsimferopole?start=ref_%d", callback.Message.Chat.ID)
	qrPath, err := h.referralService.GenerateQRCode(link, callback.Message.Chat.ID)
	if err != nil {
		h.sendError(callback.Message.Chat.ID, "Ошибка создания QR-кода.")
		return
	}

	photo := tgbotapi.NewPhoto(callback.Message.Chat.ID, tgbotapi.FilePath(qrPath))
	photo.Caption = "📷 Ваш реферальный QR-код!\nПриглашайте друзей и получайте 500 рублей за заказ от 10,000 рублей! 🎉"
	photo.ReplyMarkup = h.menus.ReferralMenu()
	if _, err := h.bot.Send(photo); err != nil {
		utils.LogError(err)
	}
}

// handleReferralPayout requests referral payout
func (h *ReferralsHandler) handleReferralPayout(callback *tgbotapi.CallbackQuery, user *models.User, data string) {
	if !h.security.HasRole(callback.Message.Chat.ID, "user") {
		h.sendUnauthorized(callback.Message.Chat.ID)
		return
	}

	inviteeIDStr := strings.TrimPrefix(data, "referral_payout_")
	inviteeID, err := strconv.Atoi(inviteeIDStr)
	if err != nil {
		h.sendError(callback.Message.Chat.ID, "Неверный формат реферала.")
		return
	}

	if err := h.referralService.RequestPayout(user.ChatID, int64(inviteeID)); err != nil {
		h.sendError(callback.Message.Chat.ID, "Ошибка запроса выплаты.")
		return
	}

	// Notify main operator and owner
	notifyMsg := fmt.Sprintf(
		"> 💸 Пользователь %s (Chat ID: %d) запрашивает выплату 500 рублей за реферал (Chat ID: %d).\n> Свяжитесь для выплаты.",
		user.FirstName, user.ChatID, inviteeID,
	)
	for _, role := range []string{"main_operator", "owner"} {
		users, _ := h.userService.ListUsersByRole(role)
		for _, u := range users {
			msg := tgbotapi.NewMessage(u.ChatID, notifyMsg)
			msg.ParseMode = "Markdown"
			h.bot.Send(msg)
		}
	}

	reply := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		"💸 Запрос на выплату отправлен. Мы свяжемся с вами! 🎉",
	)
	reply.ReplyMarkup = h.menus.ReferralMenu()
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// sendMainMenu sends the main menu
func (h *ReferralsHandler) sendMainMenu(chatID int64, user *models.User) {
	reply := tgbotapi.NewMessage(chatID, "Выберите действие:")
	reply.ReplyMarkup = h.menus.MainMenu(user)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// sendError sends an error message
func (h *ReferralsHandler) sendError(chatID int64, text string) {
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
func (h *ReferralsHandler) sendUnauthorized(chatID int64) {
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