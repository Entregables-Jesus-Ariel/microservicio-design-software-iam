package port

import (
	"context"
	"time"

	"iam/internal/domain"
)

// PasswordResetRepository handles the storage and retrieval of password reset requests.
type PasswordResetRepository interface {
	Store(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
	FindByTokenHash(ctx context.Context, tokenHash string) (domain.PasswordResetRequest, error)
	MarkAsUsed(ctx context.Context, id string) error
}
