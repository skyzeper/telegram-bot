package user

import (
	"context"
	"database/sql"
	"fmt"

	"bot/internal/models"
	"bot/internal/state"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	repo   *Repository
	states *state.StateManager
}

func NewService(db *sql.DB, states *state.StateManager) *Service {
	return &Service{
		repo:   NewRepository(db),
		states: states,
	}
}

func (s *Service) Register(ctx context.Context, chatID int64, firstName, lastName string) error {
	user := &models.User{
		ChatID:    chatID,
		Role:      "user",
		FirstName: firstName,
		LastName:  lastName,
	}
	return s.repo.CreateUser(ctx, user)
}

func (s *Service) GetUser(ctx context.Context, chatID int64) (*models.User, error) {
	return s.repo.GetUserByChatID(ctx, chatID)
}

func (s *Service) AddStaff(ctx context.Context, chatID int64, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	currentState := s.states.Get(chatID)
	switch currentState.Step {
	case 0:
		s.states.Set(chatID, state.State{
			Step:       1,
			TotalSteps: 5,
			Module:     "add_staff",
			Data:       map[string]interface{}{},
		})
		bot.Send(tgbotapi.NewMessage(chatID, "Шаг 1/5: Введите имя сотрудника"))
	case 1:
		currentState.Data["first_name"] = msg.Text
		s.states.Set(chatID, state.State{
			Step:       2,
			TotalSteps: 5,
			Module:     "add_staff",
			Data:       currentState.Data,
		})
		bot.Send(tgbotapi.NewMessage(chatID, "Шаг 2/5: Введите фамилию"))
	case 2:
		currentState.Data["last_name"] = msg.Text
		s.states.Set(chatID, state.State{
			Step:       3,
			TotalSteps: 5,
			Module:     "add_staff",
			Data:       currentState.Data,
		})
		bot.Send(tgbotapi.NewMessage(chatID, "Шаг 3/5: Введите уникальный позывной"))
	case 3:
		currentState.Data["nickname"] = msg.Text
		s.states.Set(chatID, state.State{
			Step:       4,
			TotalSteps: 5,
			Module:     "add_staff",
			Data:       currentState.Data,
		})
		bot.Send(tgbotapi.NewMessage(chatID, "Шаг 4/5: Введите телефон (+7XXX-XXX-XX-XX)"))
	case 4:
		currentState.Data["phone"] = msg.Text
		s.states.Set(chatID, state.State{
			Step:       5,
			TotalSteps: 5,
			Module:     "add_staff",
			Data:       currentState.Data,
		})
		msgConfig := tgbotapi.NewMessage(chatID, "Шаг 5/5: Выберите роль")
		msgConfig.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Оператор"),
				tgbotapi.NewKeyboardButton("Водитель"),
				tgbotapi.NewKeyboardButton("Грузчик"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Главный оператор"),
			),
		)
		bot.Send(msgConfig)
	case 5:
		role := map[string]string{
			"Оператор":         "operator",
			"Водитель":         "driver",
			"Грузчик":          "loader",
			"Главный оператор": "main_operator",
		}[msg.Text]
		user := &models.User{
			ChatID:    chatID,
			Role:      role,
			FirstName: currentState.Data["first_name"].(string),
			LastName:  currentState.Data["last_name"].(string),
			Nickname:  currentState.Data["nickname"].(string),
			Phone:     currentState.Data["phone"].(string),
		}
		if err := s.repo.CreateUser(ctx, user); err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка: %v", err)))
			return
		}
		s.states.Clear(chatID)
		msgConfig := tgbotapi.NewMessage(chatID, "Сотрудник добавлен")
		msgConfig.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{}
		bot.Send(msgConfig)
	}
}

func (s *Service) BlockUser(ctx context.Context, chatID int64, targetChatID int64, reason string) error {
	user, err := s.repo.GetUserByChatID(ctx, targetChatID)
	if err != nil {
		return err
	}
	user.IsBlocked = true
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return err
	}
	return s.repo.LogBlock(ctx, targetChatID, reason)
}
