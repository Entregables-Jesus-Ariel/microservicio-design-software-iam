// Package domain holds pure business entities and rules. It must not
// import transport, persistence or driver packages.
package domain

import "errors"

// Validation errors raised by entity constructors and behaviour.
var (
	ErrInvalidAmount        = errors.New("amount must be greater than zero")
	ErrInvalidDate          = errors.New("movement date is required")
	ErrCategoryRequired     = errors.New("movement requires a category")
	ErrUserRequired         = errors.New("movement requires a user")
	ErrNoteTooLong          = errors.New("note exceeds the maximum length of 500 characters")
	ErrInvalidCategoryType  = errors.New("category type must be income or expense")
	ErrCategoryNameRequired = errors.New("category name is required")
	ErrUsernameRequired     = errors.New("username is required")
	ErrPasswordHashRequired = errors.New("password hash is required")
	ErrInvalidRole          = errors.New("role is not supported")
	ErrInvalidAction        = errors.New("audit action must be create, update or delete")
	ErrMovementCancelled    = errors.New("movement is already cancelled")
	ErrInvalidPeriod        = errors.New("period end must not be before period start")
)

// Lookup errors raised when a persisted record does not exist.
var (
	ErrMovementNotFound   = errors.New("movement not found")
	ErrCategoryNotFound   = errors.New("category not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
