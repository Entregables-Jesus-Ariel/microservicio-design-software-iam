package domain

import "time"

// Role represents a system role in RBAC.
type Role struct {
	ID           string
	Name         string
	DisplayName  string
	Description  *string
	IsSystemRole bool
	CreatedAt    time.Time
}

// UserWithRoles is a user combined with their assigned roles.
type UserWithRoles struct {
	User
	Roles []Role
}
