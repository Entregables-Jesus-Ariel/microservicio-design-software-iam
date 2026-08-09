package usecase

import (
	"context"

	"ferreteria/internal/application/port"
	"ferreteria/internal/domain"
)

// GetMovementAudit returns the change history of a movement, oldest first.
// It is used to explain a movement when a balance does not add up.
type GetMovementAudit struct {
	movements port.MovementRepository
	audits    port.MovementAuditRepository
}

// NewGetMovementAudit wires the use case with its ports.
func NewGetMovementAudit(
	movements port.MovementRepository,
	audits port.MovementAuditRepository,
) *GetMovementAudit {
	return &GetMovementAudit{movements: movements, audits: audits}
}

// Execute validates the movement exists and returns its audit trail.
func (u *GetMovementAudit) Execute(ctx context.Context, movementID int64) ([]*domain.MovementAudit, error) {
	if _, err := u.movements.FindByID(ctx, movementID); err != nil {
		return nil, err
	}
	return u.audits.ListByMovement(ctx, movementID)
}