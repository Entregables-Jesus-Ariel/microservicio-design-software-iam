package postgres

import (
	"context"
	"database/sql"

	"iam/internal/application/port"
	"iam/internal/domain"
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

// ListLoginAttempts retrieves paginated login records and the total count.
func (r *AuditRepository) ListLoginAttempts(ctx context.Context, limit, offset int) ([]domain.LoginAttemptLog, int, error) {
	var total int
	const countQuery = `SELECT COUNT(*) FROM identity_audit.audit_login`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	const query = `
		SELECT id, user_id, email_attempted, outcome, ip_address, user_agent, attempted_at
		FROM identity_audit.audit_login
		ORDER BY attempted_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []domain.LoginAttemptLog
	for rows.Next() {
		var log domain.LoginAttemptLog
		if err := rows.Scan(
			&log.ID, &log.UserID, &log.EmailAttempted, &log.Outcome,
			&log.IPAddress, &log.UserAgent, &log.AttemptedAt,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}
	return logs, total, rows.Err()
}
