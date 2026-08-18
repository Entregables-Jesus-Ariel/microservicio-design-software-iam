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

// ListRoles returns all available roles in the system.
func (r *RoleRepository) ListRoles(ctx context.Context) ([]domain.Role, error) {
	const query = `
		SELECT id, name, display_name, description, is_system_role, created_at
		FROM rbac.role
		ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(
			&role.ID, &role.Name, &role.DisplayName, &role.Description,
			&role.IsSystemRole, &role.CreatedAt,
		); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

// RemoveRole removes a global-scope grant from rbac.user_role.
func (r *RoleRepository) RemoveRole(ctx context.Context, userID string, roleID string) error {
	const query = `
		DELETE FROM rbac.user_role
		WHERE user_id = $1 AND role_id = $2 AND training_center_id IS NULL`
	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	return err
}

// GetUserRoles returns the names of all roles assigned to a user.
func (r *RoleRepository) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	const query = `
		SELECT r.name
		FROM rbac.user_role ur
		JOIN rbac.role r ON ur.role_id = r.id
		WHERE ur.user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roleNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		roleNames = append(roleNames, name)
	}
	return roleNames, rows.Err()
}
