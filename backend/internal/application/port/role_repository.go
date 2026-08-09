package port

import "context"

// RoleRepository resolves role ids and assigns roles to users.
type RoleRepository interface {
	FindRoleIDByName(ctx context.Context, name string) (string, error)
	AssignRole(ctx context.Context, userID, roleID, assignedBy string) error
}
