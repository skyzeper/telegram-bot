package order

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/state"
	"github.com/skyzeper/telegram-bot/internal/utils"
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

func (s *Service) StartOrder(ctx context.Context, chatID int64, category string) {
	totalSteps := 9
	if category == "demolition" {
		totalSteps = 7
	}
	s.states.Set(chatID, state.State{
		Step:       1,
		TotalSteps: totalSteps,
		Module:     "create_order",
		Data: map[string]interface{}{
			"category": category,
		},
	})
}

func (s *Service) HandleOrderSteps(ctx context.Context, chatID int64, msg *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	currentState := s.states.Get(chatID)
	if currentState.Module != "create_order" {
		return
	}

	data := currentState.Data
	switch currentState.Step {
	case 1:
		data["subcategory"] = msg.Text
		s.states.Set(chatID, state.State{
			Step:       2,
			TotalSteps: currentState.TotalSteps,
			Module:     "create_order",
			Data:       data,
		})
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Шаг 2/%d: Введите ваше имя", currentState.TotalSteps)))
	case 2:
		data["name"] = msg.Text
		s.states.Set(chatID, state.State{
			Step:       3,
			TotalSteps: currentState.TotalSteps,
			Module:     "create_order",
			Data:       data,
		})
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Шаг 3/%d: Загрузите фото или видео (рекомендуем видео для точной оценки стоимости)", currentState.TotalSteps)))
	case 3:
		if msg.Photo != nil {
			photos, ok := data["photos"].([]string)
			if !ok {
				photos = []string{}
			}
			if len(photos) >= 20 {
				bot.Send(tgbotapi.NewMessage(chatID, "Достигнут лимит в 20 фото"))
				return
			}
			photos = append(photos, msg.Photo[len(msg.Photo)-1].FileID)
			data["photos"] = photos
			s.states.Set(chatID, state.State{
				Step:       4,
				TotalSteps: currentState.TotalSteps,
				Module:     "create_order",
				Data:       data,
			})
			msgConfig := tgbotapi.NewMessage(chatID, fmt.Sprintf("Шаг 4/%d: Загружено %d фото. Добавить ещё?", currentState.TotalSteps, len(photos)))
			msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("➕ Добавить", "photo_add"),
					tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить", "photo_confirm"),
					tgbotapi.NewInlineKeyboardButtonData("🔙 Отменить", "photo_cancel"),
				),
			)
			bot.Send(msgConfig)
		} else if msg.Video != nil {
			data["video"] = msg.Video.FileID
			s.states.Set(chatID, state.State{
				Step:       5,
				TotalSteps: currentState.TotalSteps,
				Module:     "create_order",
				Data:       data,
			})
			msgConfig := tgbotapi.NewMessage(chatID, fmt.Sprintf("Шаг 5/%d: Выберите дату", currentState.TotalSteps))
			msgConfig.ReplyMarkup = s.DateKeyboard()
			bot.Send(msgConfig)
		}
	case 4:
		// Обработка inline-кнопок в callbacks.go
	case 5:
		if msg.Text == "🚨 В ближайшее время" {
			data["date"] = time.Now().Format("2006-01-02")
			data["time"] = ""
			s.states.Set(chatID, state.State{
				Step:       7,
				TotalSteps: currentState.TotalSteps,
				Module:     "create_order",
				Data:       data,
			})
			msgConfig := tgbotapi.NewMessage(chatID, fmt.Sprintf("Шаг %d/%d: Введите телефон (+7XXX-XXX-XX-XX или нажмите 📞)", 7, currentState.TotalSteps))
			msgConfig.ReplyMarkup = s.phoneKeyboard()
			bot.Send(msgConfig)
		} else {
			date, err := time.Parse("2 января 2006 г.", msg.Text)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, "Неверный формат даты. Пример: 1 января 2025 г."))
				return
			}
			data["date"] = date.Format("2006-01-02")
			s.states.Set(chatID, state.State{
				Step:       6,
				TotalSteps: currentState.TotalSteps,
				Module:     "create_order",
				Data:       data,
			})
			msgConfig := tgbotapi.NewMessage(chatID, fmt.Sprintf("Шаг 6/%d: Выберите время", currentState.TotalSteps))
			msgConfig.ReplyMarkup = s.timeKeyboard()
			bot.Send(msgConfig)
		}
	case 6:
		data["time"] = msg.Text
		s.states.Set(chatID, state.State{
			Step:       7,
			TotalSteps: currentState.TotalSteps,
			Module:     "create_order",
			Data:       data,
		})
		msgConfig := tgbotapi.NewMessage(chatID, fmt.Sprintf("Шаг %d/%d: Введите телефон (+7XXX-XXX-XX-XX или нажмите 📞)", 7, currentState.TotalSteps))
		msgConfig.ReplyMarkup = s.phoneKeyboard()
		bot.Send(msgConfig)
	case 7:
		phone := utils.FormatPhone(msg.Text)
		if phone == "" && msg.Contact == nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Неверный формат телефона. Попробуйте снова"))
			return
		}
		if msg.Contact != nil {
			phone = utils.FormatPhone(msg.Contact.PhoneNumber)
		}
		data["phone"] = phone
		s.states.Set(chatID, state.State{
			Step:       8,
			TotalSteps: currentState.TotalSteps,
			Module:     "create_order",
			Data:       data,
		})
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Шаг %d/%d: Введите адрес", 8, currentState.TotalSteps)))
	case 8:
		data["address"] = msg.Text
		s.states.Set(chatID, state.State{
			Step:       9,
			TotalSteps: currentState.TotalSteps,
			Module:     "create_order",
			Data:       data,
		})
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Шаг %d/%d: Введите описание", 9, currentState.TotalSteps)))
	case 9:
		data["description"] = msg.Text
		photos, _ := data["photos"].([]string)
		video, _ := data["video"].(string)
		order := &models.Order{
			UserID:      chatID,
			Category:    data["category"].(string),
			Subcategory: data["subcategory"].(string),
			Photos:      photos,
			Video:       video,
			Date:        time.Now(),
			Time:        data["time"].(string),
			Phone:       data["phone"].(string),
			Address:     data["address"].(string),
			Description: data["description"].(string),
			Status:      "new",
		}
		if err := s.repo.CreateOrder(ctx, order); err != nil {
			msgConfig := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка сервера: %v", err))
			msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "action_back"),
				),
			)
			bot.Send(msgConfig)
			s.states.Clear(chatID)
			return
		}
		s.states.Clear(chatID)
		bot.Send(tgbotapi.NewMessage(chatID, "Заказ создан!"))
	}
}

func (s *Service) DateKeyboard() tgbotapi.InlineKeyboardMarkup {
	var buttons []tgbotapi.InlineKeyboardButton
	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("🚨 В ближайшее время", "date_urgent"))
	now := time.Now()
	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, i)
		prefix := "🟢"
		if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
			prefix = "🔴"
		}
		text := fmt.Sprintf("%s %s", prefix, date.Format("2 января 2006 г."))
		data := fmt.Sprintf("date_%s", date.Format("2006-01-02"))
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(text, data))
	}
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
}

func (s *Service) timeKeyboard() tgbotapi.InlineKeyboardMarkup {
	var buttons []tgbotapi.InlineKeyboardButton
	for hour := 9; hour <= 18; hour++ {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%02d:00", hour), fmt.Sprintf("time_%02d:00", hour)))
	}
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
}

func (s *Service) phoneKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📞 Отправить номер", "phone_contact"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "action_back"),
		),
	)
}
