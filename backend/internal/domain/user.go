package domain

import (
	"strings"
	"time"
)

// RoleAdmin is the only role supported by the current scope.
const RoleAdmin = "admin"

// User is an account allowed to record and review movements.
type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser validates and builds a user account. The password hash is
// produced outside the domain so no hashing library leaks into it.
func NewUser(username string, passwordHash string, role string) (*User, error) {
	trimmedUsername := strings.TrimSpace(username)
	if trimmedUsername == "" {
		return nil, ErrUsernameRequired
	}
	if strings.TrimSpace(passwordHash) == "" {
		return nil, ErrPasswordHashRequired
	}
	if role != RoleAdmin {
		return nil, ErrInvalidRole
	}
	now := time.Now().UTC()
	return &User{
		Username:     trimmedUsername,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// IsAdmin reports whether the user holds administrative rights.
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
