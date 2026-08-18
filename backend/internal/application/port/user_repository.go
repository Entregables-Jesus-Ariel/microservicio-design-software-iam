package port

import (
	"context"
	"time"

	"iam/internal/domain"
)

// UserRepository is the persistence contract the application layer needs
// for identity.user. Implementations live in internal/persistence.
type UserRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByID(ctx context.Context, userID string) (domain.User, error)
	RegisterFailedAttempt(ctx context.Context, userID string, lockedUntil *time.Time) error
	ResetFailedAttempts(ctx context.Context, userID string) error
	UpdatePassword(ctx context.Context, userID string, passwordHash string) error
	ListAll(ctx context.Context) ([]domain.User, error)
}
