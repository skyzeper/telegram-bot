package order_test

import (
	"errors"
	"testing"
	"time"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/services/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of order.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateOrder(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockRepository) GetOrder(id int) (*models.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockRepository) GetOrdersByStatus(status string) ([]models.Order, error) {
	args := m.Called(status)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockRepository) GetOrdersByStatusAndCategory(status, category string) ([]models.Order, error) {
	args := m.Called(status, category)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockRepository) GetExecutorOrders(userID int64) ([]models.Order, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockRepository) GetOrderClientID(orderID int) (int64, error) {
	args := m.Called(orderID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) UpdateOrder(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockRepository) ConfirmOrder(orderID int, userID int64) error {
	args := m.Called(orderID, userID)
	return args.Error(0)
}

func TestService_CreateOrder(t *testing.T) {
	mockRepo := new(MockRepository)
	service := order.NewService(mockRepo)

	t.Run("ValidOrder", func(t *testing.T) {
		order := &models.Order{
			UserID:      123,
			Category:    "вывоз мусора",
			Subcategory: "строительный мусор",
			Phone:       "+1234567890",
			Address:     "ул. Тестовая, 1",
		}

		mockRepo.On("CreateOrder", order).Return(nil).Once()

		err := service.CreateOrder(order)
		assert.NoError(t, err)
		assert.Equal(t, "new", order.Status)
		assert.False(t, order.CreatedAt.IsZero())
		assert.False(t, order.UpdatedAt.IsZero())
		mockRepo.AssertExpectations(t)
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {
		order := &models.Order{
			UserID: 123,
			Phone:  "+1234567890",
		}

		err := service.CreateOrder(order)
		assert.Error(t, err)
		assert.Equal(t, "missing required order fields", err.Error())
	})
}

func TestService_GetOrder(t *testing.T) {
	mockRepo := new(MockRepository)
	service := order.NewService(mockRepo)

	t.Run("ValidID", func(t *testing.T) {
		expectedOrder := &models.Order{
			ID:     1,
			UserID: 123,
		}

		mockRepo.On("GetOrder", 1).Return(expectedOrder, nil).Once()

		result, err := service.GetOrder(1)
		assert.NoError(t, err)
		assert.Equal(t, expectedOrder, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		result, err := service.GetOrder(0)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid order ID", err.Error())
	})
}

func TestService_UpdateOrder(t *testing.T) {
	mockRepo := new(MockRepository)
	service := order.NewService(mockRepo)

	t.Run("ValidOrder", func(t *testing.T) {
		order := &models.Order{
			ID:          1,
			UserID:      123,
			Category:    "вывоз мусора",
			Subcategory: "строительный мусор",
		}

		mockRepo.On("UpdateOrder", order).Return(nil).Once()

		originalUpdatedAt := order.UpdatedAt
		err := service.UpdateOrder(order)
		assert.NoError(t, err)
		assert.True(t, order.UpdatedAt.After(originalUpdatedAt))
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		order := &models.Order{
			ID: 0,
		}

		err := service.UpdateOrder(order)
		assert.Error(t, err)
		assert.Equal(t, "invalid order ID", err.Error())
	})
}

func TestService_GetOrdersByStatus(t *testing.T) {
	mockRepo := new(MockRepository)
	service := order.NewService(mockRepo)

	t.Run("ValidStatus", func(t *testing.T) {
		expectedOrders := []models.Order{
			{ID: 1, Status: "new"},
			{ID: 2, Status: "new"},
		}

		mockRepo.On("GetOrdersByStatus", "new").Return(expectedOrders, nil).Once()

		result, err := service.GetOrdersByStatus("new")
		assert.NoError(t, err)
		assert.Equal(t, expectedOrders, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("EmptyStatus", func(t *testing.T) {
		result, err := service.GetOrdersByStatus("")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "status cannot be empty", err.Error())
	})
}

func TestService_ConfirmOrder(t *testing.T) {
	mockRepo := new(MockRepository)
	service := order.NewService(mockRepo)

	t.Run("ValidOrder", func(t *testing.T) {
		mockRepo.On("ConfirmOrder", 1, int64(123)).Return(nil).Once()

		err := service.ConfirmOrder(1, 123)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidOrderID", func(t *testing.T) {
		err := service.ConfirmOrder(0, 123)
		assert.Error(t, err)
		assert.Equal(t, "invalid order or user ID", err.Error())
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		err := service.ConfirmOrder(1, 0)
		assert.Error(t, err)
		assert.Equal(t, "invalid order or user ID", err.Error())
	})
}