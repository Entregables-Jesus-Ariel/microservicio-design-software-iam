package port

import (
	"context"

	"ferreteria/internal/domain"
)

// UserRepository queries user accounts for authentication and audit.
type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByID(ctx context.Context, id int64) (*domain.User, error)
}
