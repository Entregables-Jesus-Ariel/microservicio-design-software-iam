package domain

import (
	"strings"
	"time"
)

// CategoryType distinguishes income from expense classification.
type CategoryType string

// Supported category types.
const (
	CategoryIncome  CategoryType = "income"
	CategoryExpense CategoryType = "expense"
)

// Valid reports whether the type is one of the supported values.
func (t CategoryType) Valid() bool {
	return t == CategoryIncome || t == CategoryExpense
}

// Category classifies a movement as income or expense.
type Category struct {
	ID        int64
	Name      string
	Type      CategoryType
	CreatedAt time.Time
}

// NewCategory validates and builds a category.
func NewCategory(name string, categoryType CategoryType) (*Category, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil, ErrCategoryNameRequired
	}
	if !categoryType.Valid() {
		return nil, ErrInvalidCategoryType
	}
	return &Category{
		Name:      trimmed,
		Type:      categoryType,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// IsIncome reports whether movements in this category add to the balance.
func (c *Category) IsIncome() bool {
	return c.Type == CategoryIncome
}
