package models

import "time"

// Executor represents an executor assigned to an order
type Executor struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	UserID    int64     `json:"user_id"`
	Role      string    `json:"role"`
	Confirmed bool      `json:"confirmed"`
	Notified  bool      `json:"notified"`
	CreatedAt time.Time `json:"created_at"`
}