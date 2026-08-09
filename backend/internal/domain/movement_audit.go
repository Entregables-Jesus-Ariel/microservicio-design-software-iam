package domain

import "time"

// AuditAction names the operation that produced an audit entry.
type AuditAction string

// Supported audit actions.
const (
	AuditCreate AuditAction = "create"
	AuditUpdate AuditAction = "update"
	AuditDelete AuditAction = "delete"
)

// Valid reports whether the action is one of the supported values.
func (a AuditAction) Valid() bool {
	return a == AuditCreate || a == AuditUpdate || a == AuditDelete
}

// MovementAudit is an immutable record of a change to a movement. It
// provides non-repudiation: every create, update and cancellation is
// traceable to a user and a timestamp.
type MovementAudit struct {
	ID         int64
	MovementID int64
	ChangedAt  time.Time
	ChangedBy  int64
	OldAmount  *Amount
	NewAmount  *Amount
	OldNote    *string
	NewNote    *string
	Action     AuditAction
}

// NewMovementAudit validates and builds an audit entry.
func NewMovementAudit(movementID int64, changedBy int64, action AuditAction) (*MovementAudit, error) {
	if movementID <= 0 {
		return nil, ErrMovementNotFound
	}
	if changedBy <= 0 {
		return nil, ErrUserRequired
	}
	if !action.Valid() {
		return nil, ErrInvalidAction
	}
	return &MovementAudit{
		MovementID: movementID,
		ChangedBy:  changedBy,
		ChangedAt:  time.Now().UTC(),
		Action:     action,
	}, nil
}

// WithAmountChange records the amount before and after the change.
func (a *MovementAudit) WithAmountChange(oldAmount Amount, newAmount Amount) *MovementAudit {
	previous := oldAmount
	current := newAmount
	a.OldAmount = &previous
	a.NewAmount = &current
	return a
}

// WithNoteChange records the note before and after the change.
func (a *MovementAudit) WithNoteChange(oldNote string, newNote string) *MovementAudit {
	previous := oldNote
	current := newNote
	a.OldNote = &previous
	a.NewNote = &current
	return a
}
