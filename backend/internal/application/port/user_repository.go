package port

import (
	"context"

	"iam/internal/domain"
)

// UserRepository is the persistence contract the application layer needs
// for identity.user. Implementations live in internal/persistence.
type UserRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}
