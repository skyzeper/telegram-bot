package stats

import (
	"fmt"
	"time"
	"github.com/skyzeper/telegram-bot/internal/models"
)

// Service handles statistics-related business logic
type Service struct {
	repo Repository
}

// Repository defines the interface for stats data access
type Repository interface {
	GetOrderStats(start, end time.Time) ([]models.Order, error)
	GetAccountingStats(start, end time.Time) ([]models.AccountingRecord, error)
}

// NewService creates a new stats service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetStatsForDay retrieves statistics for the current day
func (s *Service) GetStatsForDay() (models.Stats, error) {
	start := time.Now().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)
	return s.getStats(start, end)
}

// GetStatsForWeek retrieves statistics for the current week
func (s *Service) GetStatsForWeek() (models.Stats, error) {
	start := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -int(time.Now().Weekday()))
	end := start.AddDate(0, 0, 7)
	return s.getStats(start, end)
}

// GetStatsForMonth retrieves statistics for the current month
func (s *Service) GetStatsForMonth() (models.Stats, error) {
	start := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -time.Now().Day()+1)
	end := start.AddDate(0, 1, 0)
	return s.getStats(start, end)
}

// GetStatsForYear retrieves statistics for the current year
func (s *Service) GetStatsForYear() (models.Stats, error) {
	start := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -time.Now().Day()+1).AddDate(0, -int(time.Now().Month())+1, 0)
	end := start.AddDate(1, 0, 0)
	return s.getStats(start, end)
}

// GetStatsForAllTime retrieves statistics for all time
func (s *Service) GetStatsForAllTime() (models.Stats, error) {
	start := time.Time{} // Zero time for all history
	end := time.Now().Add(24 * time.Hour)
	return s.getStats(start, end)
}

// getStats calculates statistics for a given time range
func (s *Service) getStats(start, end time.Time) (models.Stats, error) {
	var stats models.Stats

	// Get orders
	orders, err := s.repo.GetOrderStats(start, end)
	if err != nil {
		return stats, fmt.Errorf("failed to get order stats: %v", err)
	}

	// Process orders
	for _, order := range orders {
		if order.Status == "completed" {
			stats.TotalOrders++
			switch order.Category {
			case "waste_removal":
				stats.WasteRemovalOrders++
			case "demolition":
				stats.DemolitionOrders++
			case "construction_materials":
				stats.ConstructionOrders++
			}
			stats.TotalAmount += order.Cost
		}
	}

	// Get accounting records for driver debts
	records, err := s.repo.GetAccountingStats(start, end)
	if err != nil {
		return stats, fmt.Errorf("failed to get accounting stats: %v", err)
	}

	// Process driver debts
	for _, record := range records {
		if (record.Type == "expense" || record.Type == "salary") && record.UserID != 0 {
			stats.DriverDebts += record.Amount // Negative amounts increase debt
		}
	}

	return stats, nil
}