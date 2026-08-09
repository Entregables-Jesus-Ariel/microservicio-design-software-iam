// Package usecase holds the application's business workflows. Each use
// case orchestrates ports; none of them import persistence or transport.
package usecase

import (
	"context"
	"strings"

	"iam/internal/application/port"
	"iam/internal/domain"
)

// defaultRoleName is granted automatically on public self-registration.
const defaultRoleName = "LEARNER"

// RegisterUser creates a new identity.user row and grants the default
// role, so a freshly registered account can log in immediately.
type RegisterUser struct {
	users  port.UserRepository
	roles  port.RoleRepository
	hasher port.PasswordHasher
}

// NewRegisterUser builds the use case with its dependencies.
func NewRegisterUser(users port.UserRepository, roles port.RoleRepository, hasher port.PasswordHasher) *RegisterUser {
	return &RegisterUser{users: users, roles: roles, hasher: hasher}
}

// RegisterUserInput carries the data collected from the registration form.
type RegisterUserInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

// Execute validates, persists and assigns the default role to a new user.
func (uc *RegisterUser) Execute(ctx context.Context, input RegisterUserInput) (domain.User, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	exists, err := uc.users.ExistsByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	if exists {
		return domain.User{}, domain.ErrEmailAlreadyRegistered
	}

	hash, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return domain.User{}, err
	}

	created, err := uc.users.Create(ctx, domain.User{
		Email:        email,
		PasswordHash: hash,
		FirstName:    strings.TrimSpace(input.FirstName),
		LastName:     strings.TrimSpace(input.LastName),
		ActorType:    domain.ActorTypeUser,
		IsActive:     true,
	})
	if err != nil {
		return domain.User{}, err
	}

	roleID, err := uc.roles.FindRoleIDByName(ctx, defaultRoleName)
	if err != nil {
		return domain.User{}, err
	}
	// Self-registration: the new user is its own "assigned_by" actor.
	if err := uc.roles.AssignRole(ctx, created.ID, roleID, created.ID); err != nil {
		return domain.User{}, err
	}

	return created, nil
}
