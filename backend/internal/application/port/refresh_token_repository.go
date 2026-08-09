package port

import (
	"context"
	"time"
)

// RefreshTokenRepository persists opaque refresh tokens in session.refresh_token.
type RefreshTokenRepository interface {
	Store(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
}
