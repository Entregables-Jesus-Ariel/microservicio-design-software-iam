// Package domain contains the IAM business entities. No layer below this
// one (persistence, transport) may be imported here.
package domain

import "time"

// ActorType mirrors the identity.user.actor_type CHECK constraint.
type ActorType string

const (
	ActorTypeUser       ActorType = "USER"
	ActorTypeInstructor ActorType = "INSTRUCTOR"
	ActorTypeLearner    ActorType = "LEARNER"
)

// User represents a row in identity.user.
type User struct {
	ID             string
	Email          string
	PasswordHash   string
	FirstName      string
	LastName       string
	ActorType      ActorType
	ActorID        *string
	IsActive       bool
	LastLoginAt    *time.Time
	FailedAttempts int16
	LockedUntil    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
