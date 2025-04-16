package models

import "time"

// Notification represents a notification sent to a user
type Notification struct {
	ID        int       `json:"id"`
	UserID    int64     `json:"user_id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	SentAt    time.Time `json:"sent_at"`
	CreatedAt time.Time `json:"created_at"`
}