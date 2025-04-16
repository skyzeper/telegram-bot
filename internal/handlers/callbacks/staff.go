package callbacks

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/menus"
	"github.com/skyzeper/telegram-bot/internal/services/user"
	"github.com/skyzeper/telegram-bot/internal/state"
	"github.com/skyzeper/telegram-bot/internal/security"
	"github.com/skyzeper/telegram-bot/internal/utils"
)

// StaffHandler handles staff-related callbacks
type StaffHandler struct {
	bot         *tgbotapi.BotAPI
	security    security.SecurityChecker
	menus       *menus.MenuGenerator
	userService *user.Service
	state       *state.State
}

// NewStaffHandler creates a new StaffHandler
func NewStaffHandler(
	bot *tgbotapi.BotAPI,
	security security.SecurityChecker,
	menus *menus.MenuGenerator,
	userService *user.Service,
	state *state.State,
) *StaffHandler {
	return &StaffHandler{
		bot:         bot,
		security:    security,
		menus:       menus,
		userService: userService,
		state:       state,
	}
}

// HandleStaffCallback processes staff-related callback queries
func (h *StaffHandler) HandleStaffCallback(callback *tgbotapi.CallbackQuery, user *models.User, data string) {
	switch data {
	case "staff_add":
		h.handleStaffAdd(callback, user)
	case "staff_delete":
		h.handleStaffDelete(callback, user)
	case "staff_list":
		h.handleStaffList(callback, user)
	default:
		if strings.HasPrefix(data, "staff_list_") {
			h.handleStaffListByRole(callback, user, data)
		} else if strings.HasPrefix(data, "staff_select_") {
			h.handleStaffSelect(callback, user, data)
		} else if strings.HasPrefix(data, "staff_action_") {
			h.handleStaffAction(callback, user, data)
		} else if strings.HasPrefix(data, "edit_field_") {
			h.handleEditStaffField(callback, user, data)
		} else if data == "cancel_staff" {
			h.state.ClearState(callback.Message.Chat.ID)
			h.sendMainMenu(callback.Message.Chat.ID, user)
		}
	}
}

// handleStaffAdd initiates staff addition
func (h *StaffHandler) handleStaffAdd(callback *tgbotapi.CallbackQuery, user *models.User) {
	if !h.security.HasRole(callback.Message.Chat.ID, "main_operator") && !h.security.HasRole(callback.Message.Chat.ID, "owner") {
		h.sendUnauthorized(callback.Message.Chat.ID)
		return
	}

	h.state.SetState(callback.Message.Chat.ID, state.State{
		Module:      "user",
		Step:        1,
		TotalSteps:  5,
		Data:        map[string]interface{}{"action": "add_staff"},
	})

	reply := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "Шаг 1/5: Введите имя сотрудника:")
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back_to_staff"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", "cancel_staff"),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// handleStaffDelete shows staff deletion menu
func (h *StaffHandler) handleStaffDelete(callback *tgbotapi.CallbackQuery, user *models.User) {
	if !h.security.HasRole(callback.Message.Chat.ID, "main_operator") && !h.security.HasRole(callback.Message.Chat.ID, "owner") {
		h.sendUnauthorized(callback.Message.Chat.ID)
		return
	}

	roles := []string{"operator", "driver", "loader"}
	if h.security.HasRole(callback.Message.Chat.ID, "owner") {
		roles = append(roles, "main_operator")
	}

	var buttons []tgbotapi.InlineKeyboardButton
	for _, role := range roles {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(role, fmt.Sprintf("staff_list_%s", role)))
	}

	reply := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "🗑️ Выберите роль для удаления:")
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttons...),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back_to_staff"),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// handleStaffList shows staff list
func (h *StaffHandler) handleStaffList(callback *tgbotapi.CallbackQuery, user *models.User) {
	if !h.security.HasRole(callback.Message.Chat.ID, "owner") {
		h.sendUnauthorized(callback.Message.Chat.ID)
		return
	}

	roles := []string{"operator", "main_operator", "driver", "loader"}
	var buttons []tgbotapi.InlineKeyboardButton
	for _, role := range roles {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(role, fmt.Sprintf("staff_list_%s", role)))
	}

	reply := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "📋 Выберите роль для просмотра:")
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttons...),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back_to_staff"),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// handleStaffListByRole shows staff by role
func (h *StaffHandler) handleStaffListByRole(callback *tgbotapi.CallbackQuery, user *models.User, data string) {
	if !h.security.HasRole(callback.Message.Chat.ID, "main_operator") && !h.security.HasRole(callback.Message.Chat.ID, "owner") {
		h.sendUnauthorized(callback.Message.Chat.ID)
		return
	}

	role := strings.TrimPrefix(data, "staff_list_")
	users, err := h.userService.ListUsersByRole(role)
	if err != nil {
		h.sendError(callback.Message.Chat.ID, "Ошибка получения сотрудников.")
		return
	}

	var buttons []tgbotapi.InlineKeyboardButton
	for _, u := range users {
		btnText := fmt.Sprintf("%s %s", u.FirstName, u.LastName)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(btnText, fmt.Sprintf("staff_select_%s_%d", role, u.ChatID)))
	}

	reply := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, fmt.Sprintf("📋 Сотрудники (%s):", role))
	if len(buttons) == 0 {
		reply.Text = fmt.Sprintf("📋 Нет сотрудников с ролью %s.", role)
	} else {
		reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
	}
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back_to_staff"),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// handleStaffSelect shows staff member actions
func (h *StaffHandler) handleStaffSelect(callback *tgbotapi.CallbackQuery, user *models.User, data string) {
	parts := strings.Split(strings.TrimPrefix(data, "staff_select_"), "_")
	if len(parts) != 2 {
		h.sendError(callback.Message.Chat.ID, "Неверный формат выбора.")
		return
	}

	userID, err := strconv.Atoi(parts[1])
	if err != nil {
		h.sendError(callback.Message.Chat.ID, "Неверный формат пользователя.")
		return
	}

	staff, err := h.userService.GetUser(int64(userID))
	if err != nil {
		h.sendError(callback.Message.Chat.ID, "Ошибка получения сотрудника.")
		return
	}

	reply := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		fmt.Sprintf("🧑‍💼 Сотрудник: %s %s\nРоль: %s", staff.FirstName, staff.LastName, staff.Role),
	)
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", fmt.Sprintf("staff_action_delete_%d", userID)),
			tgbotapi.NewInlineKeyboardButtonData("🚫 Заблокировать", fmt.Sprintf("staff_action_block_%d", userID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✏️ Изменить", fmt.Sprintf("staff_action_edit_%d", userID)),
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back_to_staff_list"),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// handleStaffAction performs staff actions
func (h *StaffHandler) handleStaffAction(callback *tgbotapi.CallbackQuery, user *models.User, data string) {
	parts := strings.Split(strings.TrimPrefix(data, "staff_action_"), "_")
	if len(parts) != 2 {
		h.sendError(callback.Message.Chat.ID, "Неверный формат действия.")
		return
	}

	action := parts[0]
	userID, err := strconv.Atoi(parts[1])
	if err != nil {
		h.sendError(callback.Message.Chat.ID, "Неверный формат пользователя.")
		return
	}

	switch action {
	case "delete":
		if err := h.userService.DeleteUser(int64(userID)); err != nil {
			h.sendError(callback.Message.Chat.ID, "Ошибка удаления сотрудника.")
			return
		}
		reply := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "🗑️ Сотрудник удалён.")
		reply.ReplyMarkup = h.menus.StaffMenu(callback.Message.Chat.ID)
		h.bot.Send(reply)
	case "block":
		h.state.SetState(callback.Message.Chat.ID, state.State{
			Module:      "user",
			Step:        1,
			TotalSteps:  2,
			Data:        map[string]interface{}{"action": "block_staff", "user_id": userID},
		})
		reply := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "🚫 Введите причину блокировки сотрудника:")
		reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", fmt.Sprintf("back_to_staff_select_%d", userID)),
			),
		)
		h.bot.Send(reply)
	case "edit":
		h.state.SetState(callback.Message.Chat.ID, state.State{
			Module:      "user",
			Step:        1,
			TotalSteps:  1,
			Data:        map[string]interface{}{"action": "edit_staff", "user_id": userID},
		})
		reply := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "✏️ Выберите поле для изменения:")
		reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Имя", fmt.Sprintf("edit_field_name_%d", userID)),
				tgbotapi.NewInlineKeyboardButtonData("Фамилия", fmt.Sprintf("edit_field_lastname_%d", userID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Телефон", fmt.Sprintf("edit_field_phone_%d", userID)),
				tgbotapi.NewInlineKeyboardButtonData("Роль", fmt.Sprintf("edit_field_role_%d", userID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", fmt.Sprintf("back_to_staff_select_%d", userID)),
			),
		)
		h.bot.Send(reply)
	}
}

// handleEditStaffField initiates editing a staff field
func (h *StaffHandler) handleEditStaffField(callback *tgbotapi.CallbackQuery, user *models.User, data string) {
	parts := strings.Split(strings.TrimPrefix(data, "edit_field_"), "_")
	if len(parts) != 2 {
		h.sendError(callback.Message.Chat.ID, "Неверный формат поля.")
		return
	}

	field := parts[0]
	userID, err := strconv.Atoi(parts[1])
	if err != nil {
		h.sendError(callback.Message.Chat.ID, "Неверный формат пользователя.")
		return
	}

	var fieldText string
	switch field {
	case "name":
		fieldText = "имя"
	case "lastname":
		fieldText = "фамилию"
	case "phone":
		fieldText = "телефон"
	case "role":
		fieldText = "роль (operator, driver, loader, main_operator)"
	default:
		h.sendError(callback.Message.Chat.ID, "Неверное поле.")
		return
	}

	h.state.SetState(callback.Message.Chat.ID, state.State{
		Module:      "user",
		Step:        2,
		TotalSteps:  2,
		Data:        map[string]interface{}{"action": "edit_staff", "user_id": userID, "field": field},
	})

	reply := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		fmt.Sprintf("✏️ Введите новое %s:", fieldText),
	)
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", fmt.Sprintf("back_to_edit_%d", userID)),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// sendMainMenu sends the main menu
func (h *StaffHandler) sendMainMenu(chatID int64, user *models.User) {
	reply := tgbotapi.NewMessage(chatID, "Выберите действие:")
	reply.ReplyMarkup = h.menus.MainMenu(user)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// sendError sends an error message
func (h *StaffHandler) sendError(chatID int64, text string) {
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
func (h *StaffHandler) sendUnauthorized(chatID int64) {
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