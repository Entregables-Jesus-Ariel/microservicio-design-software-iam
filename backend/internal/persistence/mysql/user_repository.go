package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ferreteria/internal/domain"
)

// UserRepository reads user accounts from MySQL.
type UserRepository struct {
	database *sql.DB
}

func NewUserRepository(database *sql.DB) *UserRepository {
	return &UserRepository{database: database}
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	const query = `
		SELECT id, username, password_hash, role, created_at, updated_at
		FROM users
		WHERE username = ?`
	return r.queryOne(ctx, query, username)
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	const query = `
		SELECT id, username, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = ?`
	return r.queryOne(ctx, query, id)
}

func (r *UserRepository) queryOne(
	ctx context.Context,
	query string,
	argument any,
) (*domain.User, error) {
	var user domain.User
	err := r.database.QueryRowContext(ctx, query, argument).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("select user: %w", err)
	}
	return &user, nil
}