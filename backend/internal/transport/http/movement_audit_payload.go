package http

import (
	"ferreteria/internal/domain"
)

// movementAuditResponse is the audit entry shape returned to clients.
// Old/new fields are nil when that side of the change does not apply
// (ej. a create entry has no old_amount or old_note).
type movementAuditResponse struct {
	ID             int64   `json:"id"`
	MovementID     int64   `json:"movement_id"`
	ChangedAt      string  `json:"changed_at"`
	ChangedBy      int64   `json:"changed_by"`
	Action         string  `json:"action"`
	OldAmountCents *int64  `json:"old_amount_cents,omitempty"`
	NewAmountCents *int64  `json:"new_amount_cents,omitempty"`
	OldNote        *string `json:"old_note,omitempty"`
	NewNote        *string `json:"new_note,omitempty"`
}

func newMovementAuditResponse(entry *domain.MovementAudit) movementAuditResponse {
	response := movementAuditResponse{
		ID:         entry.ID,
		MovementID: entry.MovementID,
		ChangedAt:  entry.ChangedAt.Format("2006-01-02T15:04:05Z07:00"),
		ChangedBy:  entry.ChangedBy,
		Action:     string(entry.Action),
		OldNote:    entry.OldNote,
		NewNote:    entry.NewNote,
	}
	if entry.OldAmount != nil {
		cents := entry.OldAmount.Cents()
		response.OldAmountCents = &cents
	}
	if entry.NewAmount != nil {
		cents := entry.NewAmount.Cents()
		response.NewAmountCents = &cents
	}
	return response
}

func newMovementAuditResponses(entries []*domain.MovementAudit) []movementAuditResponse {
	responses := make([]movementAuditResponse, 0, len(entries))
	for _, entry := range entries {
		responses = append(responses, newMovementAuditResponse(entry))
	}
	return responses
}