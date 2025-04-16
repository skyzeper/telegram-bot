package handlers_test

import (
	"errors"
	"testing"
	"time"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/handlers"
	"github.com/skyzeper/telegram-bot/internal/handlers/callbacks"
	"github.com/skyzeper/telegram-bot/internal/menus"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/security"
	"github.com/skyzeper/telegram-bot/internal/services/chat"
	"github.com/skyzeper/telegram-bot/internal/services/notification"
	"github.com/skyzeper/telegram-bot/internal/services/order"
	"github.com/skyzeper/telegram-bot/internal/services/user"
	"github.com/skyzeper/telegram-bot/internal/state"
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

// MockSecurityChecker is a mock for security.SecurityChecker
type MockSecurityChecker struct {
	mock.Mock
}

func (m *MockSecurityChecker) HasAccess(chatID int64, module string) (bool, error) {
	args := m.Called(chatID, module)
	return args.Bool(0), args.Error(1)
}

func (m *MockSecurityChecker) IsBlocked(chatID int64) (bool, error) {
	args := m.Called(chatID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSecurityChecker) GetUserRole(chatID int64) (string, error) {
	args := m.Called(chatID)
	return args.String(0), args.Error(1)
}

func (m *MockSecurityChecker) CreateUser(chatID int64, role, firstName, lastName, nickname, phone string) error {
	args := m.Called(chatID, role, firstName, lastName, nickname, phone)
	return args.Error(0)
}

// MockUserService is a mock for user.Service
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetUser(chatID int64) (*models.User, error) {
	args := m.Called(chatID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByRole(chatID int64, role string) (*models.User, error) {
	args := m.Called(chatID, role)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) ListUsersByRole(role string) ([]models.User, error) {
	args := m.Called(role)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(chatID int64) error {
	args := m.Called(chatID)
	return args.Error(0)
}

// MockChatService is a mock for chat.Service
type MockChatService struct {
	mock.Mock
}

func (m *MockChatService) SaveMessage(chatID int64, message string, fromUser bool) error {
	args := m.Called(chatID, message, fromUser)
	return args.Error(0)
}

func (m *MockChatService) GetActiveOperator() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// MockNotificationService is a mock for notification.Service
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) CreateNotification(notification *models.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func (m *MockNotificationService) ProcessPendingNotifications() error {
	args := m.Called()
	return args.Error(0)
}

func setupHandler() (*handlers.Handler, *MockBot, *MockSecurityChecker, *MockUserService, *MockChatService, *MockNotificationService) {
	mockBot := new(MockBot)
	mockSecurity := new(MockSecurityChecker)
	mockUserService := new(MockUserService)
	mockChatService := new(MockChatService)
	mockNotificationService := new(MockNotificationService)
	menus := menus.NewMenuGenerator()
	stateManager := state.NewManager()
	mockOrderService := order.NewService(new(MockRepository))
	mockCallbackHandler := &callbacks.CallbackHandler{}

	handler := handlers.NewHandler(
		mockBot,
		mockSecurity,
		menus,
		mockUserService,
		mockOrderService,
		mockChatService,
		stateManager,
		mockCallbackHandler,
		mockNotificationService,
	)
	return handler, mockBot, mockSecurity, mockUserService, mockChatService, mockNotificationService
}

func TestHandler_HandleCommand(t *testing.T) {
	handler, mockBot, mockSecurity, mockUserService, _, _ := setupHandler()

	t.Run("StartCommand", func(t *testing.T) {
		update := &tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat:      &tgbotapi.Chat{ID: 123},
				Command:   func() string { return "start" },
				IsCommand: func() bool { return true },
				From:      &tgbotapi.User{FirstName: "John", LastName: "Doe", UserName: "johndoe"},
			},
		}

		mockUserService.On("GetUser", int64(123)).Return((*models.User)(nil), errors.New("not found")).Once()
		mockSecurity.On("CreateUser", int64(123), "client", "John", "Doe", "johndoe", "").Return(nil).Once()
		mockBot.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil).Once()

		handler.HandleUpdate(update)

		mockUserService.AssertExpectations(t)
		mockSecurity.AssertExpectations(t)
		mockBot.AssertExpectations(t)
	})

	t.Run("UnknownCommand", func(t *testing.T) {
		update := &tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat:      &tgbotapi.Chat{ID: 123},
				Command:   func() string { return "unknown" },
				IsCommand: func() bool { return true },
			},
		}

		mockUserService.On("GetUser", int64(123)).Return(&models.User{ChatID: 123, Role: "client"}, nil).Once()
		mockBot.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil).Once()

		handler.HandleUpdate(update)

		mockUserService.AssertExpectations(t)
		mockBot.AssertExpectations(t)
	})
}

func TestHandler_HandleTextMessage(t *testing.T) {
	handler, mockBot, mockSecurity, mockUserService, _, _ := setupHandler()

	t.Run("OrderServiceCommand", func(t *testing.T) {
		update := &tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123},
				Text: "🗑️ заказать услугу",
			},
		}

		mockUserService.On("GetUser", int64(123)).Return(&models.User{ChatID: 123, Role: "client"}, nil).Once()
		mockSecurity.On("GetUserRole", int64(123)).Return("client", nil).Once()
		mockBot.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil).Once()

		handler.HandleUpdate(update)

		currentState := handler.State.Get(123)
		assert.Equal(t, "order", currentState.Module)
		assert.Equal(t, 1, currentState.Step)
		mockUserService.AssertExpectations(t)
		mockSecurity.AssertExpectations(t)
		mockBot.AssertExpectations(t)
	})

	t.Run("ContactOperatorCommand", func(t *testing.T) {
		update := &tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123},
				Text: "📞 связаться с оператором",
			},
		}

		mockUserService.On("GetUser", int64(123)).Return(&models.User{ChatID: 123, Role: "client"}, nil).Once()
		mockSecurity.On("GetUserRole", int64(123)).Return("client", nil).Once()
		mockBot.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil).Once()

		handler.HandleUpdate(update)

		currentState := handler.State.Get(123)
		assert.Equal(t, "chat", currentState.Module)
		assert.Equal(t, 1, currentState.Step)
		mockUserService.AssertExpectations(t)
		mockSecurity.AssertExpectations(t)
		mockBot.AssertExpectations(t)
	})

	t.Run("NoAccess", func(t *testing.T) {
		update := &tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123},
				Text: "📋 заказы",
			},
		}

		mockUserService.On("GetUser", int64(123)).Return(&models.User{ChatID: 123, Role: "client"}, nil).Once()
		mockSecurity.On("GetUserRole", int64(123)).Return("client", nil).Once()
		mockBot.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil).Once()

		handler.HandleUpdate(update)

		currentState := handler.State.Get(123)
		assert.Empty(t, currentState.Module, "state should not change")
		mockUserService.AssertExpectations(t)
		mockSecurity.AssertExpectations(t)
		mockBot.AssertExpectations(t)
	})
}

func TestHandler_HandleChatMessage(t *testing.T) {
	handler, mockBot, _, mockUserService, mockChatService, mockNotificationService := setupHandler()

	t.Run("ValidChatMessage", func(t *testing.T) {
		update := &tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123},
				Text: "Привет, нужен вывоз мусора!",
			},
		}

		handler.State.Set(123, state.State{
			Module:     "chat",
			Step:       1,
			TotalSteps: 1,
		})

		mockUserService.On("GetUser", int64(123)).Return(&models.User{ChatID: 123, Role: "client"}, nil).Once()
		mockChatService.On("SaveMessage", int64(123), "Привет, нужен вывоз мусора!", true).Return(nil).Once()
		mockChatService.On("GetActiveOperator").Return(int64(456), nil).Once()
		mockNotificationService.On("CreateNotification", mock.Anything).Return(nil).Once()
		mockBot.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil).Once()

		handler.HandleUpdate(update)

		mockUserService.AssertExpectations(t)
		mockChatService.AssertExpectations(t)
		mockNotificationService.AssertExpectations(t)
		mockBot.AssertExpectations(t)
	})

	t.Run("NoOperatorAvailable", func(t *testing.T) {
		update := &tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123},
				Text: "Привет, нужен вывоз мусора!",
			},
		}

		handler.State.Set(123, state.State{
			Module:     "chat",
			Step:       1,
			TotalSteps: 1,
		})

		mockUserService.On("GetUser", int64(123)).Return(&models.User{ChatID: 123, Role: "client"}, nil).Once()
		mockChatService.On("SaveMessage", int64(123), "Привет, нужен вывоз мусора!", true).Return(nil).Once()
		mockChatService.On("GetActiveOperator").Return(int64(0), errors.New("no operator")).Once()
		mockBot.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil).Once()

		handler.HandleUpdate(update)

		mockUserService.AssertExpectations(t)
		mockChatService.AssertExpectations(t)
		mockBot.AssertExpectations(t)
	})
}