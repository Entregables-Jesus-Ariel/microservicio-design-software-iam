package postgres

import (
	"context"
	"database/sql"
	"time"
	"errors"

	"iam/internal/domain"
)

// UserRepository is the Postgres adapter for identity.user.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository builds the repository over an open connection pool.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// ExistsByEmail reports whether a user with that email is already stored.
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM identity."user" WHERE email = $1)`
	var exists bool
	if err := r.db.QueryRowContext(ctx, query, email).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// Create inserts a new user and returns it with generated fields filled.
func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	const query = `
		INSERT INTO identity."user"
			(email, password_hash, first_name, last_name, actor_type, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	row := r.db.QueryRowContext(ctx, query,
		user.Email, user.PasswordHash, user.FirstName, user.LastName, user.ActorType, user.IsActive,
	)
	if err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// FindByEmail loads a user by email, used later by the login use case.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	const query = `
		SELECT id, email, password_hash, first_name, last_name, actor_type,
		       is_active, failed_attempts, locked_until, created_at, updated_at
		FROM identity."user"
		WHERE email = $1`

	var user domain.User
	row := r.db.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.ActorType, &user.IsActive, &user.FailedAttempts, &user.LockedUntil,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrInvalidCredentials
	}
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// FindByID loads a user by ID.
func (r *UserRepository) FindByID(ctx context.Context, userID string) (domain.User, error) {
	const query = `
		SELECT id, email, password_hash, first_name, last_name, actor_type,
		       is_active, failed_attempts, locked_until, created_at, updated_at
		FROM identity."user"
		WHERE id = $1`

	var user domain.User
	row := r.db.QueryRowContext(ctx, query, userID)
	err := row.Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.ActorType, &user.IsActive, &user.FailedAttempts, &user.LockedUntil,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrInvalidCredentials
	}
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// RegisterFailedAttempt increments failed_attempts by one and, when the
// caller decided a lockout applies, sets locked_until in the same statement.
func (r *UserRepository) RegisterFailedAttempt(ctx context.Context, userID string, lockedUntil *time.Time) error {
	const query = `
		UPDATE identity."user"
		SET failed_attempts = failed_attempts + 1,
		    locked_until = $2,
		    updated_at = now()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, userID, lockedUntil)
	return err
}

// ResetFailedAttempts clears the failure counter and any lock after a
// successful login.
func (r *UserRepository) ResetFailedAttempts(ctx context.Context, userID string) error {
	const query = `
		UPDATE identity."user"
		SET failed_attempts = 0,
		    locked_until = NULL,
		    updated_at = now()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// UpdatePassword sets a new hash for the user and unlocks them if needed.
func (r *UserRepository) UpdatePassword(ctx context.Context, userID string, passwordHash string) error {
	const query = `
		UPDATE identity."user"
		SET password_hash = $2,
		    failed_attempts = 0,
		    locked_until = NULL,
		    updated_at = now()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, userID, passwordHash)
	return err
}

// ListAll returns all users in the system.
func (r *UserRepository) ListAll(ctx context.Context) ([]domain.User, error) {
	const query = `
		SELECT id, email, password_hash, first_name, last_name, actor_type,
		       is_active, failed_attempts, locked_until, created_at, updated_at
		FROM identity."user"
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
			&user.ActorType, &user.IsActive, &user.FailedAttempts, &user.LockedUntil,
			&user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
