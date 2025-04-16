package accounting

import (
	"errors"
	"fmt"
	"time"
	"github.com/skyzeper/telegram-bot/internal/models"
)

// Service handles accounting-related business logic
type Service struct {
	repo Repository
}

// Repository defines the interface for accounting data access
type Repository interface {
	CreateRecord(record *models.AccountingRecord) error
	GetRecordsByOrder(orderID int) ([]models.AccountingRecord, error)
	GetRecordsByUser(userID int64) ([]models.AccountingRecord, error)
	GetRecordsByType(recordType string) ([]models.AccountingRecord, error)
	UpdateRecord(record *models.AccountingRecord) error
}

// NewService creates a new accounting service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateRecord creates a new accounting record
func (s *Service) CreateRecord(record *models.AccountingRecord) error {
	if record.UserID <= 0 || record.Type == "" || record.Amount == 0 {
		return errors.New("missing required accounting fields")
	}

	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now()
	}

	return s.repo.CreateRecord(record)
}

// RecordIncome records an income entry
func (s *Service) RecordIncome(orderID int, userID int64, amount float64, description string) error {
	if orderID <= 0 || userID <= 0 || amount <= 0 {
		return errors.New("invalid order ID, user ID, or amount")
	}

	record := &models.AccountingRecord{
		OrderID:     orderID,
		UserID:      userID,
		Type:        "income",
		Amount:      amount,
		Description: description,
		CreatedAt:   time.Now(),
	}

	return s.repo.CreateRecord(record)
}

// RecordExpense records an expense entry
func (s *Service) RecordExpense(userID int64, amount float64, description string) error {
	if userID <= 0 || amount <= 0 {
		return errors.New("invalid user ID or amount")
	}

	record := &models.AccountingRecord{
		UserID:      userID,
		Type:        "expense",
		Amount:      -amount, // Expenses are negative
		Description: description,
		CreatedAt:   time.Now(),
	}

	return s.repo.CreateRecord(record)
}

// RecordSalary records a salary payment
func (s *Service) RecordSalary(userID int64, amount float64, description string) error {
	if userID <= 0 || amount <= 0 {
		return errors.New("invalid user ID or amount")
	}

	record := &models.AccountingRecord{
		UserID:      userID,
		Type:        "salary",
		Amount:      -amount, // Salaries are negative (payout)
		Description: description,
		CreatedAt:   time.Now(),
	}

	return s.repo.CreateRecord(record)
}

// GetDriverDebt calculates the total debt for a driver
func (s *Service) GetDriverDebt(userID int64) (float64, error) {
	if userID <= 0 {
		return 0, errors.New("invalid user ID")
	}

	records, err := s.repo.GetRecordsByUser(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get records: %v", err)
	}

	var debt float64
	for _, record := range records {
		if record.Type == "expense" || record.Type == "salary" {
			debt += record.Amount // Negative amounts increase debt
		}
	}
	return debt, nil
}

// GetRecordsByOrder retrieves accounting records for an order
func (s *Service) GetRecordsByOrder(orderID int) ([]models.AccountingRecord, error) {
	if orderID <= 0 {
		return nil, errors.New("invalid order ID")
	}
	return s.repo.GetRecordsByOrder(orderID)
}

// GetRecordsByUser retrieves accounting records for a user
func (s *Service) GetRecordsByUser(userID int64) ([]models.AccountingRecord, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.repo.GetRecordsByUser(userID)
}

// GetRecordsByType retrieves accounting records by type
func (s *Service) GetRecordsByType(recordType string) ([]models.AccountingRecord, error) {
	if recordType == "" {
		return nil, errors.New("invalid record type")
	}
	return s.repo.GetRecordsByType(recordType)
}