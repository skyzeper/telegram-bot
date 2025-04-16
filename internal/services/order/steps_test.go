package order_test

import (
	"errors"
	"testing"
	"time"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/menus"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/state"
	"github.com/skyzeper/telegram-bot/internal/services/order"
	"github.com/skyzeper/telegram-bot/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBot is a mock implementation of tgbotapi.BotAPI
type MockBot struct {
	mock.Mock
}

func (m *MockBot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	args := m.Called(c)
	return args.Get(0).(tgbotapi.Message), args.Error(1)
}

func (m *MockBot) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	args := m.Called(c)
	return args.Get(0).(*tgbotapi.APIResponse), args.Error(1)
}

func TestStepHandler_CategoryStep(t *testing.T) {
	mockBot := new(MockBot)
	mockRepo := new(MockRepository)
	service := order.NewService(mockRepo)
	menus := menus.NewMenuGenerator()
	stateManager := state.NewManager()

	handler := order.NewStepHandler(mockBot, menus, service, stateManager)

	t.Run("ValidCategory", func(t *testing.T) {
		update := &tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123},
				Text: "Вывоз мусора",
			},
		}

		stateManager.Set(123, state.State{
			Module: "order",
			Step:   1,
		})

		mockBot.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil).Once()

		handler.HandleStep(update)

		currentState := stateManager.Get(123)
		assert.Equal(t, "order", currentState.Module)
		assert.Equal(t, 2, currentState.Step)
		assert.Equal(t, "вывоз мусора", currentState.Data["category"])
		mockBot.AssertExpectations(t)
	})

	t.Run("InvalidCategory", func(t *testing.T) {
		update := &tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123},
				Text: "Неправильная категория",
			},
		}

		stateManager.Set(123, state.State{
			Module: "order",
			Step:   1,
		})

		mockBot.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil).Once()

		handler.HandleStep(update)

		currentState := stateManager.Get(123)
		assert.Equal(t, "order", currentState.Module)
		assert.Equal(t, 1, currentState.Step)
		assert.Nil(t, currentState.Data["category"])
		mockBot.AssertExpectations(t)
	})
}

func TestStepHandler_ConfirmationStep(t *testing.T) {
	mockBot := new(MockBot)
	mockRepo := new(MockRepository)
	service := order.NewService(mockRepo)
	menus := menus.NewMenuGenerator()
	stateManager := state.NewManager()

	handler := order.NewStepHandler(mockBot, menus, service, stateManager)

	t.Run("ValidConfirmation", func(t *testing.T) {
		update := &tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123},
				Text: "Подтвердить",
			},
		}

		stateManager.Set(123, state.State{
			Module:     "order",
			Step:       11,
			TotalSteps: 11,
			Data: map[string]interface{}{
				"category":      "вывоз мусора",
				"subcategory":   "строительный мусор",
				"date":          "2025-04-17",
				"time":          "14:00",
				"phone":         "+1234567890",
				"address":       "ул. Тестовая, 1",
				"description":   "Тестовый заказ",
				"payment_method": "наличные",
				"photos":        []string{"photo1"},
				"video":         "",
			},
		})

		mockRepo.On("CreateOrder", mock.Anything).Return(nil).Once()
		mockBot.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil).Once()

		handler.HandleStep(update)

		currentState := stateManager.Get(123)
		assert.Empty(t, currentState.Module, "state should be cleared")
		mockRepo.AssertExpectations(t)
		mockBot.AssertExpectations(t)
	})

	t.Run("InvalidTime", func(t *testing.T) {
		update := &tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123},
				Text: "Подтвердить",
			},
		}

		stateManager.Set(123, state.State{
			Module:     "order",
			Step:       11,
			TotalSteps: 11,
			Data: map[string]interface{}{
				"category":      "вывоз мусора",
				"subcategory":   "строительный мусор",
				"date":          "2025-04-17",
				"time":          "invalid",
				"phone":         "+1234567890",
				"address":       "ул. Тестовая, 1",
				"description":   "Тестовый заказ",
				"payment_method": "наличные",
			},
		})

		user := &models.User{ChatID: 123}
		mockBot.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil).Once()

		handler.HandleStep(update)

		currentState := stateManager.Get(123)
		assert.Equal(t, "order", currentState.Module, "state should not be cleared")
		mockBot.AssertExpectations(t)
	})
}