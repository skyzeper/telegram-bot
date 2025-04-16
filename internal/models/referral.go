package models

import "time"

// Referral represents a referral record
type Referral struct {
	ID             int       `json:"id"`
	InviterID      int64     `json:"inviter_id"`
	InviteeID      int64     `json:"invitee_id"`
	OrderID        int       `json:"order_id"`
	PayoutRequested bool      `json:"payout_requested"`
	CreatedAt      time.Time `json:"created_at"`
}