package http

import (
	"time"

	"ferreteria/internal/domain"
)

// movementRequest is the accepted shape for creating or editing a movement.
// Only these fields are read, so a client cannot set identifiers or
// ownership by adding extra keys.
type movementRequest struct {
	CategoryID  int64  `json:"category_id"`
	Date        string `json:"date"`
	AmountCents int64  `json:"amount_cents"`
	Note        string `json:"note"`
}

// parsedDate converts the request date, which must be an ISO calendar day.
func (r movementRequest) parsedDate() (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		return time.Time{}, domain.ErrInvalidDate
	}
	return parsed, nil
}

// movementResponse is the movement shape returned to clients.
type movementResponse struct {
	ID          int64  `json:"id"`
	CategoryID  int64  `json:"category_id"`
	Date        string `json:"date"`
	AmountCents int64  `json:"amount_cents"`
	Note        string `json:"note"`
	Cancelled   bool   `json:"cancelled"`
}

func newMovementResponse(movement *domain.Movement) movementResponse {
	return movementResponse{
		ID:          movement.ID,
		CategoryID:  movement.CategoryID,
		Date:        movement.Date.Format("2006-01-02"),
		AmountCents: movement.Amount.Cents(),
		Note:        movement.Note,
		Cancelled:   movement.IsCancelled(),
	}
}

func newMovementResponses(movements []*domain.Movement) []movementResponse {
	responses := make([]movementResponse, 0, len(movements))
	for _, movement := range movements {
		responses = append(responses, newMovementResponse(movement))
	}
	return responses
}
