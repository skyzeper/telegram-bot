package models

import "time"

type Payment struct {
	ID        int64     `json:"id"`
	OrderID   int64     `json:"order_id"`
	UserID    int64     `json:"user_id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	DriverID  int64     `json:"driver_id"`
	Confirmed bool      `json:"confirmed"`
	CreatedAt time.Time `json:"created_at"`
}
