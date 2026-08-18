package domain

import "time"

// LoginAttemptLog represents a detailed record of a login attempt for audit UI.
type LoginAttemptLog struct {
	ID             string
	UserID         *string
	EmailAttempted string
	Outcome        string
	IPAddress      *string
	UserAgent      *string
	AttemptedAt    time.Time
}

// PaginatedLogins returns a slice of logs and the total count.
type PaginatedLogins struct {
	Items []LoginAttemptLog
	Total int
	Page  int
	Limit int
}
