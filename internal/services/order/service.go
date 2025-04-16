package order

import (
	"errors"
	"time"
	"github.com/skyzeper/telegram-bot/internal/models"
)

// Service handles order-related business logic
type Service struct {
	repo Repository
}

// Repository defines the interface for order data access
type Repository interface {
	CreateOrder(order *models.Order) error
	GetOrder(id int) (*models.Order, error)
	GetOrdersByStatus(status string) ([]models.Order, error)
	GetOrdersByStatusAndCategory(status, category string) ([]models.Order, error)
	GetExecutorOrders(userID int64) ([]models.Order, error)
	GetOrderClientID(orderID int) (int64, error)
	UpdateOrder(order *models.Order) error
	ConfirmOrder(orderID int, userID int64) error
}

// NewService creates a new order service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateOrder creates a new order
func (s *Service) CreateOrder(order *models.Order) error {
	if order.UserID == 0 || order.Category == "" || order.Subcategory == "" || order.Phone == "" || order.Address == "" {
		return errors.New("missing required order fields")
	}

	if order.Status == "" {
		order.Status = "new"
	}
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
	}
	if order.UpdatedAt.IsZero() {
		order.UpdatedAt = time.Now()
	}

	return s.repo.CreateOrder(order)
}

// GetOrder retrieves an order by ID
func (s *Service) GetOrder(id int) (*models.Order, error) {
	if id <= 0 {
		return nil, errors.New("invalid order ID")
	}
	return s.repo.GetOrder(id)
}

// GetOrdersByStatus retrieves orders by status
func (s *Service) GetOrdersByStatus(status string) ([]models.Order, error) {
	if status == "" {
		return nil, errors.New("status cannot be empty")
	}
	return s.repo.GetOrdersByStatus(status)
}

// GetOrdersByStatusAndCategory retrieves orders by status and category
func (s *Service) GetOrdersByStatusAndCategory(status, category string) ([]models.Order, error) {
	if status == "" || category == "" {
		return nil, errors.New("status and category cannot be empty")
	}
	return s.repo.GetOrdersByStatusAndCategory(status, category)
}

// GetExecutorOrders retrieves orders assigned to an executor
func (s *Service) GetExecutorOrders(userID int64) ([]models.Order, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.repo.GetExecutorOrders(userID)
}

// GetOrderClientID retrieves the client ID for an order
func (s *Service) GetOrderClientID(orderID int) (int64, error) {
	if orderID <= 0 {
		return 0, errors.New("invalid order ID")
	}
	return s.repo.GetOrderClientID(orderID)
}

// UpdateOrder updates an existing order
func (s *Service) UpdateOrder(order *models.Order) error {
	if order.ID <= 0 {
		return errors.New("invalid order ID")
	}
	order.UpdatedAt = time.Now()
	return s.repo.UpdateOrder(order)
}

// ConfirmOrder confirms an order
func (s *Service) ConfirmOrder(orderID int, userID int64) error {
	if orderID <= 0 || userID <= 0 {
		return errors.New("invalid order or user ID")
	}
	return s.repo.ConfirmOrder(orderID, userID)
}