package models

import "time"

type AccountingRecord struct {
	ID          int64     `json:"id"`
	OrderID     int64     `json:"order_id"`
	UserID      int64     `json:"user_id"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
