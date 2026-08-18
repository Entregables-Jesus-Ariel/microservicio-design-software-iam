package domain

import "time"

// PasswordResetRequest represents a row in session.password_reset_request.
type PasswordResetRequest struct {
	ID          string
	UserID      string
	TokenHash   string
	ExpiresAt   time.Time
	IsUsed      bool
	RequestedAt time.Time
	IPAddress   *string
}
