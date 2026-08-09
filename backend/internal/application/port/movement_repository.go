// Package port declares the interfaces the application layer needs. The
// concrete implementations live in the persistence and security layers.
package port

import (
	"context"

	"ferreteria/internal/domain"
)

// MovementFilter bounds a movement listing query.
type MovementFilter struct {
	Period           domain.Period
	CategoryID       int64
	IncludeCancelled bool
	Limit            int
	Offset           int
}

// CategoryTotal aggregates movement amounts for one category type.
type CategoryTotal struct {
	Type  domain.CategoryType
	Cents int64
}

// MovementRepository persists and queries movements.
type MovementRepository interface {
	Create(ctx context.Context, movement *domain.Movement) (*domain.Movement, error)
	Update(ctx context.Context, movement *domain.Movement) error
	FindByID(ctx context.Context, id int64) (*domain.Movement, error)
	List(ctx context.Context, filter MovementFilter) ([]*domain.Movement, error)
	TotalsByCategoryType(ctx context.Context, period domain.Period) ([]CategoryTotal, error)
}
