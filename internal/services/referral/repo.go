package referral

import (
	"context"
	"database/sql"

	"github.com/skyzeper/telegram-bot/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateReferral(ctx context.Context, referral *models.Referral) error {
	query := `
		INSERT INTO referrals (inviter_id, invitee_id)
		VALUES ($1, $2)
		RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query,
		referral.InviterID, referral.InviteeID,
	).Scan(&referral.ID, &referral.CreatedAt)
}
