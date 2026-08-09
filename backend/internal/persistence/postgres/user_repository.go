package postgres

import (
	"context"
	"database/sql"
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
