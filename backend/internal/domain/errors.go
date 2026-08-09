package domain

import "errors"

// Sentinel errors the application layer maps to HTTP status codes.
var (
	ErrEmailAlreadyRegistered = errors.New("email is already registered")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrAccountLocked          = errors.New("account is temporarily locked")
	ErrDefaultRoleMissing     = errors.New("default role LEARNER is not seeded")
)
