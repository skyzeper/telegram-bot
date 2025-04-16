package order

import (
	"fmt"
	"strings"
	"time"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/menus"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/state"
	"github.com/skyzeper/telegram-bot/internal/utils"
)

// StepHandler handles the step-by-step order creation process
type StepHandler struct {
	bot     *tgbotapi.BotAPI
	menus   *menus.MenuGenerator
	service *Service
	state   *state.Manager
}

// NewStepHandler creates a new StepHandler
func NewStepHandler(bot *tgbotapi.BotAPI, menus *menus.MenuGenerator, service *Service, state *state.Manager) *StepHandler {
	return &StepHandler{
		bot:     bot,
		menus:   menus,
		service: service,
		state:   state,
	}
}

// HandleStep processes the current step in the order creation process
func (h *StepHandler) HandleStep(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	currentState := h.state.Get(chatID)
	if currentState.Module != "order" {
		return
	}

	switch currentState.Step {
	case 1: // Category
		h.handleCategoryStep(update)
	case 2: // Subcategory
		h.handleSubcategoryStep(update)
	case 3: // Date
		h.handleDateStep(update)
	case 4: // Time
		h.handleTimeStep(update)
	case 5: // Photos
		h.handlePhotosStep(update)
	case 6: // Video
		h.handleVideoStep(update)
	case 7: // Phone
		h.handlePhoneStep(update)
	case 8: // Address
		h.handleAddressStep(update)
	case 9: // Description
		h.handleDescriptionStep(update)
	case 10: // Payment Method
		h.handlePaymentMethodStep(update)
	case 11: // Confirmation
		h.handleConfirmationStep(update)
	}
}

// handleCategoryStep handles the category selection step
func (h *StepHandler) handleCategoryStep(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	if update.Message.Text == "" {
		h.sendStepMessage(chatID, "🗑️ Выберите категорию заказа:", h.menus.CategoryMenu())
		return
	}

	category := strings.ToLower(update.Message.Text)
	if category != "вывоз мусора" && category != "демонтаж" && category != "стройматериалы" {
		h.sendStepMessage(chatID, "❌ Неверная категория. Выберите из предложенных:", h.menus.CategoryMenu())
		return
	}

	h.state.Set(chatID, state.State{
		Module:     "order",
		Step:       2,
		TotalSteps: 11,
		Data: map[string]interface{}{
			"category": category,
		},
	})
	h.sendStepMessage(chatID, "🔍 Выберите подкатегорию:", h.menus.SubcategoryMenu(category))
}

// handleSubcategoryStep handles the subcategory selection step
func (h *StepHandler) handleSubcategoryStep(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	currentState := h.state.Get(chatID)
	category, _ := currentState.Data["category"].(string)

	if update.Message.Text == "" {
		h.sendStepMessage(chatID, "🔍 Выберите подкатегорию:", h.menus.SubcategoryMenu(category))
		return
	}

	subcategory := strings.ToLower(update.Message.Text)
	valid := false
	switch category {
	case "вывоз мусора":
		valid = subcategory == "строительный мусор" || subcategory == "бытовой мусор" || subcategory == "мебель"
	case "демонтаж":
		valid = subcategory == "стены" || subcategory == "полы" || subcategory == "потолки"
	case "стройматериалы":
		valid = subcategory == "песок" || subcategory == "цемент" || subcategory == "кирпич"
	}

	if !valid {
		h.sendStepMessage(chatID, "❌ Неверная подкатегория. Выберите из предложенных:", h.menus.SubcategoryMenu(category))
		return
	}

	currentState.Data["subcategory"] = subcategory
	currentState.Step = 3
	h.state.Set(chatID, currentState)
	h.sendStepMessage(chatID, "📅 Выберите дату заказа:", h.menus.DateMenu())
}

// handleDateStep handles the date selection step
func (h *StepHandler) handleDateStep(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	if update.Message.Text == "" {
		h.sendStepMessage(chatID, "📅 Выберите дату заказа:", h.menus.DateMenu())
		return
	}

	date, err := time.Parse("02.01.2006", update.Message.Text)
	if err != nil || date.Before(time.Now().Truncate(24*time.Hour)) {
		h.sendStepMessage(chatID, "❌ Неверная или прошедшая дата. Выберите из предложенных:", h.menus.DateMenu())
		return
	}

	currentState := h.state.Get(chatID)
	currentState.Data["date"] = date.Format("2006-01-02")
	currentState.Step = 4
	h.state.Set(chatID, currentState)
	h.sendStepMessage(chatID, "🕒 Выберите время заказа:", h.menus.TimeMenu())
}

// handleTimeStep handles the time selection step
func (h *StepHandler) handleTimeStep(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	if update.Message.Text == "" {
		h.sendStepMessage(chatID, "🕒 Выберите время заказа:", h.menus.TimeMenu())
		return
	}

	_, err := time.Parse("15:04", update.Message.Text)
	if err != nil {
		h.sendStepMessage(chatID, "❌ Неверное время. Выберите из предложенных:", h.menus.TimeMenu())
		return
	}

	currentState := h.state.Get(chatID)
	currentState.Data["time"] = update.Message.Text
	currentState.Step = 5
	h.state.Set(chatID, currentState)
	h.sendStepMessage(chatID, "📸 Прикрепите фотографии (или пропустите):", h.menus.PhotoMenu())
}

// handlePhotosStep handles the photo upload step
func (h *StepHandler) handlePhotosStep(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	currentState := h.state.Get(chatID)

	if update.Message.Text == "Пропустить" {
		currentState.Data["photos"] = []string{}
		currentState.Step = 6
		h.state.Set(chatID, currentState)
		h.sendStepMessage(chatID, "🎥 Прикрепите видео (или пропустите):", h.menus.VideoMenu())
		return
	}

	if update.Message.Photo != nil {
		var photos []string
		if p, ok := currentState.Data["photos"].([]string); ok {
			photos = p
		}
		for _, photo := range update.Message.Photo {
			photos = append(photos, photo.FileID)
		}
		currentState.Data["photos"] = photos
		h.state.Set(chatID, currentState)
		h.sendStepMessage(chatID, "📸 Прикрепите ещё фото или пропустите:", h.menus.PhotoMenu())
		return
	}

	h.sendStepMessage(chatID, "📸 Прикрепите фотографии или пропустите:", h.menus.PhotoMenu())
}

// handleVideoStep handles the video upload step
func (h *StepHandler) handleVideoStep(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	currentState := h.state.Get(chatID)

	if update.Message.Text == "Пропустить" {
		currentState.Data["video"] = ""
		currentState.Step = 7
		h.state.Set(chatID, currentState)
		h.sendStepMessage(chatID, "📞 Введите номер телефона:", h.menus.PhoneMenu())
		return
	}

	if update.Message.Video != nil {
		currentState.Data["video"] = update.Message.Video.FileID
		currentState.Step = 7
		h.state.Set(chatID, currentState)
		h.sendStepMessage(chatID, "📞 Введите номер телефона:", h.menus.PhoneMenu())
		return
	}

	h.sendStepMessage(chatID, "🎥 Прикрепите видео или пропустите:", h.menus.VideoMenu())
}

// handlePhoneStep handles the phone number input step
func (h *StepHandler) handlePhoneStep(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	currentState := h.state.Get(chatID)

	if update.Message.Text == "" && update.Message.Contact == nil {
		h.sendStepMessage(chatID, "📞 Введите номер телефона или отправьте контакт:", h.menus.PhoneMenu())
		return
	}

	phone := update.Message.Text
	if update.Message.Contact != nil {
		phone = update.Message.Contact.PhoneNumber
	}

	if !utils.IsValidPhone(phone) {
		h.sendStepMessage(chatID, "❌ Неверный формат телефона. Введите корректный номер:", h.menus.PhoneMenu())
		return
	}

	currentState.Data["phone"] = phone
	currentState.Step = 8
	h.state.Set(chatID, currentState)
	h.sendStepMessage(chatID, "📍 Введите адрес:", nil)
}

// handleAddressStep handles the address input step
func (h *StepHandler) handleAddressStep(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	if update.Message.Text == "" {
		h.sendStepMessage(chatID, "📍 Введите адрес:", nil)
		return
	}

	currentState := h.state.Get(chatID)
	currentState.Data["address"] = update.Message.Text
	currentState.Step = 9
	h.state.Set(chatID, currentState)
	h.sendStepMessage(chatID, "💬 Введите описание заказа (или пропустите):", h.menus.SkipMenu())
}

// handleDescriptionStep handles the description input step
func (h *StepHandler) handleDescriptionStep(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	currentState := h.state.Get(chatID)

	if update.Message.Text == "Пропустить" {
		currentState.Data["description"] = ""
	} else {
		currentState.Data["description"] = update.Message.Text
	}

	currentState.Step = 10
	h.state.Set(chatID, currentState)
	h.sendStepMessage(chatID, "💳 Выберите способ оплаты:", h.menus.PaymentMenu())
}

// handlePaymentMethodStep handles the payment method selection step
func (h *StepHandler) handlePaymentMethodStep(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	if update.Message.Text == "" {
		h.sendStepMessage(chatID, "💳 Выберите способ оплаты:", h.menus.PaymentMenu())
		return
	}

	paymentMethod := strings.ToLower(update.Message.Text)
	if paymentMethod != "наличные" && paymentMethod != "карта" {
		h.sendStepMessage(chatID, "❌ Неверный способ оплаты. Выберите из предложенных:", h.menus.PaymentMenu())
		return
	}

	currentState := h.state.Get(chatID)
	currentState.Data["payment_method"] = paymentMethod
	currentState.Step = 11
	h.state.Set(chatID, currentState)
	h.handleConfirmationStep(update)
}

// handleConfirmationStep handles the final confirmation step
func (h *StepHandler) handleConfirmationStep(update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	currentState := h.state.Get(chatID)

	if update.Message.Text != "Подтвердить" {
		summary := h.generateOrderSummary(currentState.Data)
		h.sendStepMessage(chatID, fmt.Sprintf("📋 Подтвердите заказ:\n%s", summary), h.menus.ConfirmMenu())
		return
	}

	// Parse time
	timeVal, err := time.Parse("15:04", currentState.Data["time"].(string))
	if err != nil {
		user := &models.User{ChatID: chatID}
		h.sendStepMessage(chatID, "❌ Ошибка обработки времени. Попробуйте снова:", h.menus.MainMenu(user))
		return
	}

	// Create order
	date, _ := time.Parse("2006-01-02", currentState.Data["date"].(string))
	order := &models.Order{
		UserID:        chatID,
		Category:      currentState.Data["category"].(string),
		Subcategory:   currentState.Data["subcategory"].(string),
		Photos:        convertToStringSlice(currentState.Data["photos"]),
		Video:         currentState.Data["video"].(string),
		Date:          date,
		Time:          timeVal,
		Phone:         currentState.Data["phone"].(string),
		Address:       currentState.Data["address"].(string),
		Description:   currentState.Data["description"].(string),
		PaymentMethod: currentState.Data["payment_method"].(string),
	}

	if err := h.service.CreateOrder(order); err != nil {
		user := &models.User{ChatID: chatID}
		h.sendStepMessage(chatID, "❌ Ошибка создания заказа. Попробуйте снова:", h.menus.MainMenu(user))
		utils.LogError(err)
		return
	}

	h.state.Clear(chatID)
	user := &models.User{ChatID: chatID}
	reply := tgbotapi.NewMessage(chatID, "✅ Заказ успешно создан! Мы свяжемся для подтверждения стоимости. 😊")
	reply.ReplyMarkup = h.menus.MainMenu(user)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// generateOrderSummary creates a summary of the order
func (h *StepHandler) generateOrderSummary(data map[string]interface{}) string {
	return fmt.Sprintf(
		"> **Детали заказа** 🚛\n"+
			"> Категория: %s (%s)\n"+
			"> Дата: %s\n"+
			"> Время: %s\n"+
			"> Телефон: %s\n"+
			"> Адрес: %s\n"+
			"> Описание: %s\n"+
			"> Способ оплаты: %s",
		data["category"], data["subcategory"], data["date"], data["time"],
		data["phone"], data["address"], data["description"], data["payment_method"],
	)
}

// convertToStringSlice converts interface to string slice
func convertToStringSlice(i interface{}) []string {
	if i == nil {
		return nil
	}
	if slice, ok := i.([]string); ok {
		return slice
	}
	return nil
}

// sendStepMessage sends a message for the current step
func (h *StepHandler) sendStepMessage(chatID int64, text string, replyMarkup interface{}) {
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
		utils.LogError(err)
	}
}