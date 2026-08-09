package port

import (
	"context"

	"ferreteria/internal/domain"
)

// MovementAuditRepository appends immutable audit entries. It exposes no
// update or delete operation so the trail cannot be rewritten.
type MovementAuditRepository interface {
	Append(ctx context.Context, entry *domain.MovementAudit) error
	ListByMovement(ctx context.Context, movementID int64) ([]*domain.MovementAudit, error)
}
