package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"ferreteria/internal/domain"
)

// movementRecord is the persistence shape of a movement row.
type movementRecord struct {
	ID          int64
	UserID      int64
	CategoryID  int64
	Date        time.Time
	AmountCents int64
	Note        sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CancelledAt sql.NullTime
}

func (r movementRecord) toDomain() *domain.Movement {
	movement := &domain.Movement{
		ID:         r.ID,
		UserID:     r.UserID,
		CategoryID: r.CategoryID,
		Date:       r.Date,
		Amount:     domain.AmountFromCents(r.AmountCents),
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
	if r.Note.Valid {
		movement.Note = r.Note.String
	}
	if r.CancelledAt.Valid {
		cancelled := r.CancelledAt.Time
		movement.CancelledAt = &cancelled
	}
	return movement
}

func noteValue(note string) sql.NullString {
	if note == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: note, Valid: true}
}

func cancelledValue(cancelledAt *time.Time) sql.NullTime {
	if cancelledAt == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *cancelledAt, Valid: true}
}

// centsToDecimalString renders integer cents as a fixed 2-decimal string
// (ej. 123456 -> "1234.56") para que el monto llegue a la columna DECIMAL
// de MySQL como texto exacto, no como float64 con riesgo de redondeo.
func centsToDecimalString(cents int64) string {
	negative := cents < 0
	if negative {
		cents = -cents
	}
	whole := cents / 100
	fraction := cents % 100
	result := fmt.Sprintf("%d.%02d", whole, fraction)
	if negative {
		result = "-" + result
	}
	return result
}