package order

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"bot/internal/models"
	"bot/internal/state"
	"bot/internal/utils"
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

func (s *Service) StartOrder(ctx context.Context, chatID int64, category string) {
	s.states.Set(chatID, state.State{
		Step:       1,
		TotalSteps: 9,
		Module:     "create_order",
		Data: map[string]interface{}{
			"category": category,
		},
	})
}

func (s *Service) HandleOrderSteps(chatID int64, msg *tgbotapi.Message, bot *tgbotapi.BotAPI) {
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
			TotalSteps: 9,
			Module:     "create_order",
			Data:       data,
		})
		bot.Send(tgbotapi.NewMessage(chatID, "Шаг 2/9: Введите ваше имя"))
	case 2:
		data["name"] = msg.Text
		s.states.Set(chatID, state.State{
			Step:       3,
			TotalSteps: 9,
			Module:     "create_order",
			Data:       data,
		})
		bot.Send(tgbotapi.NewMessage(chatID, "Шаг 3/9: Загрузите фото или видео (рекомендуем видео для точной оценки стоимости)"))
	case 3:
		if msg.Photo != nil {
			photos, ok := data["photos"].([]string)
			if !ok {
				photos = []string{}
			}
			photos = append(photos, msg.Photo[len(msg.Photo)-1].FileID)
			data["photos"] = photos
			s.states.Set(chatID, state.State{
				Step:       4,
				TotalSteps: 9,
				Module:     "create_order",
				Data:       data,
			})
			bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Шаг 4/9: Загружено %d фото. Добавить ещё?", len(photos)),
				tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("➕ Добавить"),
						tgbotapi.NewKeyboardButton("✅ Подтвердить"),
						tgbotapi.NewKeyboardButton("🔙 Отменить"),
					),
				)))
		} else if msg.Video != nil {
			data["video"] = msg.Video.FileID
			s.states.Set(chatID, state.State{
				Step:       5,
				TotalSteps: 9,
				Module:     "create_order",
				Data:       data,
			})
			bot.Send(tgbotapi.NewMessage(chatID, "Шаг 5/9: Выберите дату", s.dateKeyboard())))
		}
	case 4:
		if msg.Text == "✅ Подтвердить" {
			s.states.Set(chatID, state.State{
				Step:       5,
				TotalSteps: 9,
				Module:     "create_order",
				Data:       data,
			})
			bot.Send(tgbotapi.NewMessage(chatID, "Шаг 5/9: Выберите дату", s.dateKeyboard()))
		} else if msg.Text == "🔙 Отменить" {
			s.states.Clear(chatID)
			bot.Send(tgbotapi.NewMessage(chatID, "Создание заказа отменено", tgbotapi.ReplyKeyboardRemove{}))
		}
	case 5:
		if msg.Text == "🚨 В ближайшее время" {
			data["date"] = time.Now().Format("2006-01-02")
			data["time"] = ""
			s.states.Set(chatID, state.State{
				Step:       7,
				TotalSteps: 9,
				Module:     "create_order",
				Data:       data,
			})
			bot.Send(tgbotapi.NewMessage(chatID, "Шаг 7/9: Введите телефон (+7XXX-XXX-XX-XX или нажмите 📞)", s.phoneKeyboard()))
		} else {
			date, err := time.Parse("2 января 2006 г.", msg.Text)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, "Неверный формат даты. Попробуйте снова"))
				return
			}
			data["date"] = date.Format("2006-01-02")
			s.states.Set(chatID, state.State{
				Step:       6,
				TotalSteps: 9,
				Module:     "create_order",
				Data:       data,
			})
			bot.Send(tgbotapi.NewMessage(chatID, "Шаг 6/9: Выберите время", s.timeKeyboard()))
		}
	case 6:
		data["time"] = msg.Text
		s.states.Set(chatID, state.State{
			Step:       7,
			TotalSteps: 9,
			Module:     "create_order",
			Data:       data,
		})
		bot.Send(tgbotapi.NewMessage(chatID, "Шаг 7/9: Введите телефон (+7XXX-XXX-XX-XX или нажмите 📞)", s.phoneKeyboard()))
	case 7:
		phone := utils.FormatPhone(msg.Text)
		if phone == "" {
			bot.Send(tgbotapi.NewMessage(chatID, "Неверный формат телефона. Попробуйте снова"))
			return
		}
		data["phone"] = phone
		s.states.Set(chatID, state.State{
			Step:       8,
			TotalSteps: 9,
			Module:     "create_order",
			Data:       data,
		})
		bot.Send(tgbotapi.NewMessage(chatID, "Шаг 8/9: Введите адрес"))
	case 8:
		data["address"] = msg.Text
		s.states.Set(chatID, state.State{
			Step:       9,
			TotalSteps: 9,
			Module:     "create_order",
			Data:       data,
		})
		bot.Send(tgbotapi.NewMessage(chatID, "Шаг 9/9: Введите описание"))
	case 9:
		data["description"] = msg.Text
		order := &models.Order{
			UserID:      chatID,
			Category:    data["category"].(string),
			Subcategory: data["subcategory"].(string),
			Photos:      data["photos"].([]string),
			Video:       data["video"].(string),
			Date:        time.Now(),
			Time:        data["time"].(string),
			Phone:       data["phone"].(string),
			Address:     data["address"].(string),
			Description: data["description"].(string),
			Status:      "new",
		}
		if err := s.repo.CreateOrder(ctx, order); err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка: %v", err)))
			return
		}
		s.states.Clear(chatID)
		bot.Send(tgbotapi.NewMessage(chatID, "Заказ создан!", tgbotapi.ReplyKeyboardRemove{}))
	}
}

func (s *Service) dateKeyboard() tgbotapi.InlineKeyboardMarkup {
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

func (s *Service) phoneKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("📞 Отправить номер"),
		),
	)
}
