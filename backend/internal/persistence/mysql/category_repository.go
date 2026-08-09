package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ferreteria/internal/domain"
)

// CategoryRepository stores categories in MySQL.
type CategoryRepository struct {
	database *sql.DB
}

func NewCategoryRepository(database *sql.DB) *CategoryRepository {
	return &CategoryRepository{database: database}
}

func (r *CategoryRepository) Create(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	const statement = `
		INSERT INTO categories (name, type, created_at)
		VALUES (?, ?, ?)`

	result, err := r.database.ExecContext(
		ctx,
		statement,
		category.Name,
		string(category.Type),
		category.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert category: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("read inserted category id: %w", err)
	}
	category.ID = id
	return category, nil
}

func (r *CategoryRepository) FindByID(ctx context.Context, id int64) (*domain.Category, error) {
	const query = `SELECT id, name, type, created_at FROM categories WHERE id = ?`
	return r.queryOne(ctx, query, id)
}

func (r *CategoryRepository) FindByName(ctx context.Context, name string) (*domain.Category, error) {
	const query = `SELECT id, name, type, created_at FROM categories WHERE name = ?`
	return r.queryOne(ctx, query, name)
}

func (r *CategoryRepository) ListByType(
	ctx context.Context,
	categoryType domain.CategoryType,
) ([]*domain.Category, error) {
	const query = `
		SELECT id, name, type, created_at
		FROM categories
		WHERE type = ?
		ORDER BY name`
	return r.queryMany(ctx, query, string(categoryType))
}

func (r *CategoryRepository) ListAll(ctx context.Context) ([]*domain.Category, error) {
	const query = `SELECT id, name, type, created_at FROM categories ORDER BY name`
	return r.queryMany(ctx, query)
}

func (r *CategoryRepository) queryOne(
	ctx context.Context,
	query string,
	argument any,
) (*domain.Category, error) {
	var category domain.Category
	var categoryType string
	err := r.database.QueryRowContext(ctx, query, argument).Scan(
		&category.ID,
		&category.Name,
		&categoryType,
		&category.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrCategoryNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("select category: %w", err)
	}
	category.Type = domain.CategoryType(categoryType)
	return &category, nil
}

func (r *CategoryRepository) queryMany(
	ctx context.Context,
	query string,
	arguments ...any,
) ([]*domain.Category, error) {
	rows, err := r.database.QueryContext(ctx, query, arguments...)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	categories := make([]*domain.Category, 0)
	for rows.Next() {
		var category domain.Category
		var categoryType string
		if err := rows.Scan(&category.ID, &category.Name, &categoryType, &category.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan category row: %w", err)
		}
		category.Type = domain.CategoryType(categoryType)
		categories = append(categories, &category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate category rows: %w", err)
	}
	return categories, nil
}