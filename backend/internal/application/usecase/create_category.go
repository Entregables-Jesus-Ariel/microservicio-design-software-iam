package usecase

import (
	"context"
	"errors"

	"ferreteria/internal/application/port"
	"ferreteria/internal/domain"
)

// ErrCategoryNameTaken is returned when the name already exists.
var ErrCategoryNameTaken = errors.New("category name is already in use")

// CreateCategory adds a custom classification.
type CreateCategory struct {
	categories port.CategoryRepository
}

// NewCreateCategory wires the use case with its ports.
func NewCreateCategory(categories port.CategoryRepository) *CreateCategory {
	return &CreateCategory{categories: categories}
}

// Execute validates uniqueness and stores the category.
func (u *CreateCategory) Execute(
	ctx context.Context,
	name string,
	categoryType domain.CategoryType,
) (*domain.Category, error) {
	category, err := domain.NewCategory(name, categoryType)
	if err != nil {
		return nil, err
	}

	existing, err := u.categories.FindByName(ctx, category.Name)
	if err != nil && !errors.Is(err, domain.ErrCategoryNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrCategoryNameTaken
	}

	return u.categories.Create(ctx, category)
}
