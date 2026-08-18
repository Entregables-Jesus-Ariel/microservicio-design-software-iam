package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"iam/internal/domain"
)

// RefreshTokenRepository is the Postgres adapter for session.refresh_token.
type RefreshTokenRepository struct {
	db *sql.DB
}

// NewRefreshTokenRepository builds the repository over an open connection pool.
func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Store persists a refresh token's hash, never the plain value.
func (r *RefreshTokenRepository) Store(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error {
	const query = `
		INSERT INTO session.refresh_token (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, userID, tokenHash, expiresAt)
	return err
}

// Find retrieves the user ID associated with a valid (not revoked, not expired) token hash.
func (r *RefreshTokenRepository) Find(ctx context.Context, tokenHash string) (string, error) {
	const query = `
		SELECT user_id 
		FROM session.refresh_token 
		WHERE token_hash = $1 AND is_revoked = FALSE AND expires_at > now()`

	var userID string
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrInvalidToken
		}
		return "", err
	}

	return userID, nil
}

// Revoke marks a token as revoked.
func (r *RefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	const query = `
		UPDATE session.refresh_token 
		SET is_revoked = TRUE, revoked_at = now() 
		WHERE token_hash = $1`

	_, err := r.db.ExecContext(ctx, query, tokenHash)
	return err
}

