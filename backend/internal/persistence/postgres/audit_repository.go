package postgres

import (
	"context"
	"database/sql"

	"iam/internal/application/port"
)

// AuditRepository is the Postgres adapter for identity_audit.audit_login.
type AuditRepository struct {
	db *sql.DB
}

// NewAuditRepository builds the repository over an open connection pool.
func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// RecordLoginAttempt writes one row per authentication attempt, successful or not.
func (r *AuditRepository) RecordLoginAttempt(ctx context.Context, attempt port.LoginAttempt) error {
	const query = `
		INSERT INTO identity_audit.audit_login (user_id, email_attempted, outcome)
		VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, attempt.UserID, attempt.EmailAttempted, attempt.Outcome)
	return err
}
