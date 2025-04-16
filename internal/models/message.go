package models

import "time"

// Message represents a chat message
type Message struct {
	ID         int       `json:"id"`
	UserID     int64     `json:"user_id"`
	OperatorID int64     `json:"operator_id"`
	Message    string    `json:"message"`
	IsFromUser bool      `json:"is_from_user"`
	CreatedAt  time.Time `json:"created_at"`
}