package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"github.com/xuri/excelize/v2"
	"github.com/skyzeper/telegram-bot/internal/models"
)

// GenerateExcelReport creates an Excel report for accounting records
func GenerateExcelReport(records []models.AccountingRecord, filename string) (string, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			LogError(err)
		}
	}()

	// Create a new sheet
	sheet := "Report"
	f.SetSheetName("Sheet1", sheet)

	// Set headers
	headers := []string{"ID", "Order ID", "User ID", "Type", "Amount", "Description", "Created At"}
	for col, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+col)))
		f.SetCellValue(sheet, cell, header)
	}

	// Fill data
	for row, record := range records {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row+2), record.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row+2), record.OrderID)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row+2), record.UserID)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row+2), record.Type)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row+2), record.Amount)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row+2), record.Description)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row+2), record.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	// Auto-size columns
	for col := 'A'; col <= 'G'; col++ {
		f.SetColWidth(sheet, string(col), string(col), 15)
	}

	// Save file
	filepath := filepath.Join(os.TempDir(), fmt.Sprintf("report_%s.xlsx", time.Now().Format("20060102150405")))
	if err := f.SaveAs(filepath); err != nil {
		return "", fmt.Errorf("failed to save Excel file: %v", err)
	}

	return filepath, nil
}