package usecase

import (
	"context"
	"time"

	"ferreteria/internal/application/port"
	"ferreteria/internal/domain"
)

// CancelMovement marks a movement as cancelled while preserving it.
type CancelMovement struct {
	movements port.MovementRepository
	audits    port.MovementAuditRepository
}

// NewCancelMovement wires the use case with its ports.
func NewCancelMovement(
	movements port.MovementRepository,
	audits port.MovementAuditRepository,
) *CancelMovement {
	return &CancelMovement{movements: movements, audits: audits}
}

// Execute cancels the movement and records the cancellation.
func (u *CancelMovement) Execute(ctx context.Context, movementID int64, actorID int64) error {
	movement, err := u.movements.FindByID(ctx, movementID)
	if err != nil {
		return err
	}

	if err := movement.Cancel(time.Now().UTC()); err != nil {
		return err
	}
	if err := u.movements.Update(ctx, movement); err != nil {
		return err
	}

	entry, err := domain.NewMovementAudit(movement.ID, actorID, domain.AuditDelete)
	if err != nil {
		return err
	}
	previous := movement.Amount
	entry.OldAmount = &previous
	return u.audits.Append(ctx, entry)
}
