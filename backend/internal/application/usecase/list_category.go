package usecase

import (
	"context"

	"ferreteria/internal/application/port"
	"ferreteria/internal/domain"
)

// ListCategory returns the categories available for classification.
type ListCategory struct {
	categories port.CategoryRepository
}

// NewListCategory wires the use case with its ports.
func NewListCategory(categories port.CategoryRepository) *ListCategory {
	return &ListCategory{categories: categories}
}

// Execute returns every category, or only those of one type when a valid
// filter is supplied.
func (u *ListCategory) Execute(
	ctx context.Context,
	filter domain.CategoryType,
) ([]*domain.Category, error) {
	if filter == "" {
		return u.categories.ListAll(ctx)
	}
	if !filter.Valid() {
		return nil, domain.ErrInvalidCategoryType
	}
	return u.categories.ListByType(ctx, filter)
}
