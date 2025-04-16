package referral

import (
	"database/sql"
	"fmt"
	"github.com/skyzeper/telegram-bot/internal/db"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/utils"
)

// PostgresRepository implements the Repository interface for PostgreSQL
type PostgresRepository struct {
	db *db.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *db.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateReferral creates a new referral
func (r *PostgresRepository) CreateReferral(referral *models.Referral) error {
	query := `
		INSERT INTO referrals (inviter_id, invitee_id, order_id, payout_requested, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.Conn().QueryRow(
		query,
		referral.InviterID, referral.InviteeID, referral.OrderID, referral.PayoutRequested, referral.CreatedAt,
	).Scan(&referral.ID)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to create referral: %v", err)
	}
	return nil
}

// GetReferralByInvitee retrieves a referral by invitee ID
func (r *PostgresRepository) GetReferralByInvitee(inviteeID int64) (*models.Referral, error) {
	query := `
		SELECT id, inviter_id, invitee_id, order_id, payout_requested, created_at
		FROM referrals
		WHERE invitee_id = $1
	`
	referral := &models.Referral{}
	err := r.db.Conn().QueryRow(query, inviteeID).Scan(
		&referral.ID, &referral.InviterID, &referral.InviteeID, &referral.OrderID,
		&referral.PayoutRequested, &referral.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("referral not found")
	}
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get referral: %v", err)
	}
	return referral, nil
}

// GetReferralsByInviter retrieves all referrals by inviter
func (r *PostgresRepository) GetReferralsByInviter(inviterID int64) ([]models.Referral, error) {
	query := `
		SELECT id, inviter_id, invitee_id, order_id, payout_requested, created_at
		FROM referrals
		WHERE inviter_id = $1
	`
	rows, err := r.db.Conn().Query(query, inviterID)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get referrals by inviter: %v", err)
	}
	defer rows.Close()

	var referrals []models.Referral
	for rows.Next() {
		var referral models.Referral
		if err := rows.Scan(
			&referral.ID, &referral.InviterID, &referral.InviteeID, &referral.OrderID,
			&referral.PayoutRequested, &referral.CreatedAt,
		); err != nil {
			utils.LogError(err)
			continue
		}
		referrals = append(referrals, referral)
	}
	return referrals, nil
}

// UpdateReferral updates an existing referral
func (r *PostgresRepository) UpdateReferral(referral *models.Referral) error {
	query := `
		UPDATE referrals
		SET inviter_id = $1, invitee_id = $2, order_id = $3, payout_requested = $4, created_at = $5
		WHERE id = $6
	`
	_, err := r.db.Conn().Exec(
		query,
		referral.InviterID, referral.InviteeID, referral.OrderID, referral.PayoutRequested,
		referral.CreatedAt, referral.ID,
	)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to update referral: %v", err)
	}
	return nil
}