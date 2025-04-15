package accounting

import (
	"bot/internal/models"
	"context"
	"database/sql"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(db *sql.DB) *Service {
	return &Service{
		repo: NewRepository(db),
	}
}

func (s *Service) LogExpense(ctx context.Context, orderID, userID int64, expenseType string, amount float64, description string) error {
	record := &models.AccountingRecord{
		OrderID:     orderID,
		UserID:      userID,
		Type:        expenseType,
		Amount:      amount,
		Description: description,
	}
	if err := s.repo.CreateRecord(ctx, record); err != nil {
		return err
	}
	return s.updateExcel(record)
}

func (s *Service) updateExcel(record *models.AccountingRecord) error {
	f := excelize.NewFile()
	defer f.Close()
	f.SetCellValue("Sheet1", "A1", "Date")
	f.SetCellValue("Sheet1", "B1", "OrderID")
	f.SetCellValue("Sheet1", "C1", "UserID")
	f.SetCellValue("Sheet1", "D1", "Type")
	f.SetCellValue("Sheet1", "E1", "Amount")
	f.SetCellValue("Sheet1", "F1", "Description")
	f.SetCellValue("Sheet1", "A2", record.CreatedAt.Format("2006-01-02"))
	f.SetCellValue("Sheet1", "B2", record.OrderID)
	f.SetCellValue("Sheet1", "C2", record.UserID)
	f.SetCellValue("Sheet1", "D2", record.Type)
	f.SetCellValue("Sheet1", "E2", record.Amount)
	f.SetCellValue("Sheet1", "F2", record.Description)
	return f.SaveAs(fmt.Sprintf("bot/accounting/drivers/%s/driver_%d.xlsx", record.CreatedAt.Format("January_2006"), record.UserID))
}
