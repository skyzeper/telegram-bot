package models

import "time"

type Executor struct {
	ID        int64     `json:"id"`
	OrderID   int64     `json:"order_id"`
	UserID    int64     `json:"user_id"`
	Role      string    `json:"role"`
	Confirmed bool      `json:"confirmed"`
	Notified  bool      `json:"notified"`
	CreatedAt time.Time `json:"created_at"`
}
