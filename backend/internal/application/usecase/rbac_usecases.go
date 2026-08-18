package usecase

import (
	"context"

	"iam/internal/application/port"
	"iam/internal/domain"
)

type RBACUsecases struct {
	users port.UserRepository
	roles port.RoleRepository
}

func NewRBACUsecases(users port.UserRepository, roles port.RoleRepository) *RBACUsecases {
	return &RBACUsecases{users: users, roles: roles}
}

// ListUsersWithRoles returns all users and their assigned roles.
func (uc *RBACUsecases) ListUsersWithRoles(ctx context.Context) ([]domain.UserWithRoles, error) {
	users, err := uc.users.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []domain.UserWithRoles
	for _, u := range users {
		roleNames, err := uc.roles.GetUserRoles(ctx, u.ID)
		if err != nil {
			return nil, err
		}
		
		var roles []domain.Role
		for _, name := range roleNames {
			roles = append(roles, domain.Role{Name: name})
		}
		
		result = append(result, domain.UserWithRoles{
			User:  u,
			Roles: roles,
		})
	}
	return result, nil
}

// ListRoles returns the catalog of available roles.
func (uc *RBACUsecases) ListRoles(ctx context.Context) ([]domain.Role, error) {
	return uc.roles.ListRoles(ctx)
}

type AssignRoleInput struct {
	UserID     string
	RoleID     string
	AssignedBy string
}

// AssignRole assigns a role to a user globally.
func (uc *RBACUsecases) AssignRole(ctx context.Context, input AssignRoleInput) error {
	return uc.roles.AssignRole(ctx, input.UserID, input.RoleID, input.AssignedBy)
}

// RevokeRole removes a global role from a user.
func (uc *RBACUsecases) RevokeRole(ctx context.Context, userID, roleID string) error {
	return uc.roles.RemoveRole(ctx, userID, roleID)
}
