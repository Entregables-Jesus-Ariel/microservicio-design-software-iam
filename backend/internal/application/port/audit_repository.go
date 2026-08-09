package port

import "context"

// LoginOutcome mirrors the identity_audit.audit_login.outcome CHECK constraint.
type LoginOutcome string

const (
	LoginOutcomeSuccess         LoginOutcome = "SUCCESS"
	LoginOutcomeInvalidPassword LoginOutcome = "INVALID_PASSWORD"
	LoginOutcomeUserNotFound    LoginOutcome = "USER_NOT_FOUND"
	LoginOutcomeAccountLocked   LoginOutcome = "ACCOUNT_LOCKED"
)

// LoginAttempt is one row written to identity_audit.audit_login.
type LoginAttempt struct {
	UserID         *string
	EmailAttempted string
	Outcome        LoginOutcome
}

// AuditRepository records authentication events.
type AuditRepository interface {
	RecordLoginAttempt(ctx context.Context, attempt LoginAttempt) error
}
