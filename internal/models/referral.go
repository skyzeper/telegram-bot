package models

import "time"

type Referral struct {
	ID              int64     `json:"id"`
	InviterID       int64     `json:"inviter_id"`
	InviteeID       int64     `json:"invitee_id"`
	OrderID         int64     `json:"order_id"`
	PayoutRequested bool      `json:"payout_requested"`
	CreatedAt       time.Time `json:"created_at"`
}
