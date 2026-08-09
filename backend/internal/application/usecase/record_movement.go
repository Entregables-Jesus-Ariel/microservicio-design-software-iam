// Package usecase holds the application rules that orchestrate domain
// entities through repository ports. It performs no HTTP or SQL work.
package usecase

import (
	"context"
	"time"

	"ferreteria/internal/application/port"
	"ferreteria/internal/domain"
)

// RecordMovementInput describes a new income or expense entry.
type RecordMovementInput struct {
	UserID       int64
	CategoryID   int64
	Date         time.Time
	AmountCents  int64
	Note         string
}

// RecordMovement stores a movement and its creation audit entry.
type RecordMovement struct {
	movements  port.MovementRepository
	categories port.CategoryRepository
	audits     port.MovementAuditRepository
}

// NewRecordMovement wires the use case with its ports.
func NewRecordMovement(
	movements port.MovementRepository,
	categories port.CategoryRepository,
	audits port.MovementAuditRepository,
) *RecordMovement {
	return &RecordMovement{movements: movements, categories: categories, audits: audits}
}

// Execute validates the input and persists the movement.
func (u *RecordMovement) Execute(ctx context.Context, input RecordMovementInput) (*domain.Movement, error) {
	if _, err := u.categories.FindByID(ctx, input.CategoryID); err != nil {
		return nil, err
	}

	amount, err := domain.NewAmount(input.AmountCents)
	if err != nil {
		return nil, err
	}

	movement, err := domain.NewMovement(input.UserID, input.CategoryID, input.Date, amount, input.Note)
	if err != nil {
		return nil, err
	}

	stored, err := u.movements.Create(ctx, movement)
	if err != nil {
		return nil, err
	}

	entry, err := domain.NewMovementAudit(stored.ID, input.UserID, domain.AuditCreate)
	if err != nil {
		return nil, err
	}
	entry.NewAmount = &stored.Amount
	note := stored.Note
	entry.NewNote = &note
	if err := u.audits.Append(ctx, entry); err != nil {
		return nil, err
	}

	return stored, nil
}
