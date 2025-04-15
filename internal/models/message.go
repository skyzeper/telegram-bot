package models

import "time"

type Message struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	OperatorID int64     `json:"operator_id"`
	Message    string    `json:"message"`
	IsFromUser bool      `json:"is_from_user"`
	CreatedAt  time.Time `json:"created_at"`
}
