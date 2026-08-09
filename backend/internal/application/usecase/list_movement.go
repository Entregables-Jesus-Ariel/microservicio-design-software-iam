package usecase

import (
	"context"
	"time"

	"ferreteria/internal/application/port"
	"ferreteria/internal/domain"
)

// DefaultPageSize bounds a listing page so a large table cannot be pulled
// in a single request.
const DefaultPageSize = 50

// ListMovementInput describes a movement listing query.
type ListMovementInput struct {
	Start      time.Time
	End        time.Time
	CategoryID int64
	Page       int
	PageSize   int
}

// ListMovement returns the movements recorded inside a period.
type ListMovement struct {
	movements port.MovementRepository
}

// NewListMovement wires the use case with its ports.
func NewListMovement(movements port.MovementRepository) *ListMovement {
	return &ListMovement{movements: movements}
}

// Execute validates the period and returns one page of movements.
func (u *ListMovement) Execute(ctx context.Context, input ListMovementInput) ([]*domain.Movement, error) {
	period, err := domain.NewPeriod(input.Start, input.End)
	if err != nil {
		return nil, err
	}

	pageSize := input.PageSize
	if pageSize <= 0 || pageSize > DefaultPageSize {
		pageSize = DefaultPageSize
	}
	page := input.Page
	if page < 1 {
		page = 1
	}

	return u.movements.List(ctx, port.MovementFilter{
		Period:     period,
		CategoryID: input.CategoryID,
		// El movimiento anulado se marca como tal pero permanece visible en
		// el listado (HU-008), así que siempre se incluye.
		IncludeCancelled: true,
		Limit:            pageSize,
		Offset:           (page - 1) * pageSize,
	})
}
