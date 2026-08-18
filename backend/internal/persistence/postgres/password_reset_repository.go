package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"iam/internal/domain"
)

// PasswordResetRepository is the Postgres adapter for session.password_reset_request.
type PasswordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

// Store persists a new password reset request.
func (r *PasswordResetRepository) Store(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error {
	const query = `
		INSERT INTO session.password_reset_request (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, userID, tokenHash, expiresAt)
	return err
}

// FindByTokenHash loads the reset request by the hashed token.
func (r *PasswordResetRepository) FindByTokenHash(ctx context.Context, tokenHash string) (domain.PasswordResetRequest, error) {
	const query = `
		SELECT id, user_id, token_hash, expires_at, is_used, requested_at, ip_address
		FROM session.password_reset_request
		WHERE token_hash = $1`

	var req domain.PasswordResetRequest
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&req.ID, &req.UserID, &req.TokenHash, &req.ExpiresAt, &req.IsUsed, &req.RequestedAt, &req.IPAddress,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PasswordResetRequest{}, domain.ErrInvalidToken
		}
		return domain.PasswordResetRequest{}, err
	}

	return req, nil
}

// MarkAsUsed sets the request as used so it cannot be used again.
func (r *PasswordResetRepository) MarkAsUsed(ctx context.Context, id string) error {
	const query = `
		UPDATE session.password_reset_request
		SET is_used = TRUE
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
