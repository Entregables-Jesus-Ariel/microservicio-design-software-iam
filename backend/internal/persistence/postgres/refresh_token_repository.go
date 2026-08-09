package postgres

import (
	"context"
	"database/sql"
	"time"
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
