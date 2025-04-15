package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	Role      string    `json:"role"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Nickname  string    `json:"nickname"`
	Phone     string    `json:"phone"`
	IsBlocked bool      `json:"is_blocked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
