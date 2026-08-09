package usecase

import (
	"context"
	"time"

	"ferreteria/internal/application/port"
	"ferreteria/internal/domain"
)

// EditMovementInput describes the new values for a movement.
type EditMovementInput struct {
	MovementID  int64
	EditorID    int64
	CategoryID  int64
	Date        time.Time
	AmountCents int64
	Note        string
}

// EditMovement updates a movement and records the change.
type EditMovement struct {
	movements  port.MovementRepository
	categories port.CategoryRepository
	audits     port.MovementAuditRepository
}

// NewEditMovement wires the use case with its ports.
func NewEditMovement(
	movements port.MovementRepository,
	categories port.CategoryRepository,
	audits port.MovementAuditRepository,
) *EditMovement {
	return &EditMovement{movements: movements, categories: categories, audits: audits}
}

// Execute applies the edit after validating the movement and category.
func (u *EditMovement) Execute(ctx context.Context, input EditMovementInput) (*domain.Movement, error) {
	movement, err := u.movements.FindByID(ctx, input.MovementID)
	if err != nil {
		return nil, err
	}
	if _, err := u.categories.FindByID(ctx, input.CategoryID); err != nil {
		return nil, err
	}

	amount, err := domain.NewAmount(input.AmountCents)
	if err != nil {
		return nil, err
	}

	previousAmount := movement.Amount
	previousNote := movement.Note

	if err := movement.Edit(input.CategoryID, input.Date, amount, input.Note); err != nil {
		return nil, err
	}
	if err := u.movements.Update(ctx, movement); err != nil {
		return nil, err
	}

	entry, err := domain.NewMovementAudit(movement.ID, input.EditorID, domain.AuditUpdate)
	if err != nil {
		return nil, err
	}
	entry.WithAmountChange(previousAmount, movement.Amount).WithNoteChange(previousNote, movement.Note)
	if err := u.audits.Append(ctx, entry); err != nil {
		return nil, err
	}

	return movement, nil
}
