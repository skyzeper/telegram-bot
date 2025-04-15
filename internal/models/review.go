package models

import "time"

type Review struct {
	ID        int64     `json:"id"`
	OrderID   int64     `json:"order_id"`
	UserID    int64     `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
