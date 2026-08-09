package postgres

import (
	"context"
	"database/sql"
	"errors"

	"iam/internal/domain"
)

// RoleRepository is the Postgres adapter for rbac.role / rbac.user_role.
type RoleRepository struct {
	db *sql.DB
}

// NewRoleRepository builds the repository over an open connection pool.
func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// FindRoleIDByName resolves a role id from its unique name (e.g. LEARNER).
func (r *RoleRepository) FindRoleIDByName(ctx context.Context, name string) (string, error) {
	const query = `SELECT id FROM rbac.role WHERE name = $1`
	var id string
	err := r.db.QueryRowContext(ctx, query, name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", domain.ErrDefaultRoleMissing
	}
	if err != nil {
		return "", err
	}
	return id, nil
}

// AssignRole inserts a global-scope grant into rbac.user_role.
func (r *RoleRepository) AssignRole(ctx context.Context, userID, roleID, assignedBy string) error {
	const query = `
		INSERT INTO rbac.user_role (user_id, role_id, assigned_by)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, role_id, training_center_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, userID, roleID, assignedBy)
	return err
}
