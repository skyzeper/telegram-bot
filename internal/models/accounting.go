package models

import "time"

// AccountingRecord represents an accounting entry
type AccountingRecord struct {
	ID          int       `json:"id"`
	OrderID     int       `json:"order_id"`
	UserID      int64     `json:"user_id"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}