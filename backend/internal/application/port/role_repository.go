package port

import (
	"context"

	"iam/internal/domain"
)

// RoleRepository resolves role ids and assigns roles to users.
type RoleRepository interface {
	FindRoleIDByName(ctx context.Context, name string) (string, error)
	ListRoles(ctx context.Context) ([]domain.Role, error)
	AssignRole(ctx context.Context, userID, roleID, assignedBy string) error
	RemoveRole(ctx context.Context, userID string, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
}
