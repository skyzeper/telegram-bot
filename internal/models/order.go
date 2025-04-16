package models

import "time"

// Order represents an order in the system
type Order struct {
	ID            int        `json:"id"`
	UserID        int64      `json:"user_id"`
	Category      string     `json:"category"`
	Subcategory   string     `json:"subcategory"`
	Photos        []string   `json:"photos"`
	Video         string     `json:"video"`
	Date          time.Time  `json:"date"`
	Time          time.Time  `json:"time"`
	Phone         string     `json:"phone"`
	Address       string     `json:"address"`
	Description   string     `json:"description"`
	Status        string     `json:"status"`
	Reason        string     `json:"reason"`
	Cost          float64    `json:"cost"`
	PaymentMethod string     `json:"payment_method"`
	PaymentConfirmed bool    `json:"payment_confirmed"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Executors     []Executor `json:"executors"`
	Confirmed     bool       `json:"confirmed"`
}