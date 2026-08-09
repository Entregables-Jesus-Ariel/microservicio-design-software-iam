package port

import (
	"context"

	"ferreteria/internal/domain"
)

// CategoryRepository persists and queries categories.
type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) (*domain.Category, error)
	FindByID(ctx context.Context, id int64) (*domain.Category, error)
	FindByName(ctx context.Context, name string) (*domain.Category, error)
	ListByType(ctx context.Context, categoryType domain.CategoryType) ([]*domain.Category, error)
	ListAll(ctx context.Context) ([]*domain.Category, error)
}
