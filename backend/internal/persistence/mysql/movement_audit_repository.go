package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"ferreteria/internal/domain"
)

// MovementAuditRepository appends audit rows. It exposes no update or
// delete operation so the trail stays immutable.
type MovementAuditRepository struct {
	database *sql.DB
}

func NewMovementAuditRepository(database *sql.DB) *MovementAuditRepository {
	return &MovementAuditRepository{database: database}
}

func (r *MovementAuditRepository) Append(ctx context.Context, entry *domain.MovementAudit) error {
	const statement = `
		INSERT INTO movement_audit (
			movement_id, changed_at, changed_by,
			old_amount, new_amount, old_note, new_note, action
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.database.ExecContext(
		ctx,
		statement,
		entry.MovementID,
		entry.ChangedAt,
		entry.ChangedBy,
		amountDecimalValue(entry.OldAmount),
		amountDecimalValue(entry.NewAmount),
		noteValue(optionalString(entry.OldNote)),
		noteValue(optionalString(entry.NewNote)),
		string(entry.Action),
	)
	if err != nil {
		return fmt.Errorf("insert movement audit: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("read inserted audit id: %w", err)
	}
	entry.ID = id
	return nil
}

func (r *MovementAuditRepository) ListByMovement(
	ctx context.Context,
	movementID int64,
) ([]*domain.MovementAudit, error) {
	const query = `
		SELECT id, movement_id, changed_at, changed_by,
			old_amount, new_amount, old_note, new_note, action
		FROM movement_audit
		WHERE movement_id = ?
		ORDER BY changed_at, id`

	rows, err := r.database.QueryContext(ctx, query, movementID)
	if err != nil {
		return nil, fmt.Errorf("list movement audit: %w", err)
	}
	defer rows.Close()

	entries := make([]*domain.MovementAudit, 0)
	for rows.Next() {
		var entry domain.MovementAudit
		var action string
		var oldAmount, newAmount sql.NullString
		var oldNote, newNote sql.NullString
		if err := rows.Scan(
			&entry.ID,
			&entry.MovementID,
			&entry.ChangedAt,
			&entry.ChangedBy,
			&oldAmount,
			&newAmount,
			&oldNote,
			&newNote,
			&action,
		); err != nil {
			return nil, fmt.Errorf("scan audit row: %w", err)
		}
		entry.Action = domain.AuditAction(action)
		if amount, ok, err := decimalStringToAmount(oldAmount); err != nil {
			return nil, fmt.Errorf("parse old_amount: %w", err)
		} else if ok {
			entry.OldAmount = &amount
		}
		if amount, ok, err := decimalStringToAmount(newAmount); err != nil {
			return nil, fmt.Errorf("parse new_amount: %w", err)
		} else if ok {
			entry.NewAmount = &amount
		}
		if oldNote.Valid {
			entry.OldNote = &oldNote.String
		}
		if newNote.Valid {
			entry.NewNote = &newNote.String
		}
		entries = append(entries, &entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate audit rows: %w", err)
	}
	return entries, nil
}

func amountDecimalValue(amount *domain.Amount) sql.NullString {
	if amount == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: centsToDecimalString(amount.Cents()), Valid: true}
}

func optionalString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// decimalStringToAmount parses a DECIMAL column such as "1234.56" back into
// a domain.Amount expressed in whole cents. It returns ok=false when the
// column was NULL, which happens for the side of the change that did not
// apply (ej. a create entry has no old_amount).
func decimalStringToAmount(value sql.NullString) (domain.Amount, bool, error) {
	if !value.Valid {
		return domain.Amount{}, false, nil
	}
	raw := value.String
	negative := strings.HasPrefix(raw, "-")
	if negative {
		raw = raw[1:]
	}
	parts := strings.SplitN(raw, ".", 2)
	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return domain.Amount{}, false, fmt.Errorf("parse whole part %q: %w", raw, err)
	}
	var fraction int64
	if len(parts) == 2 {
		fractionDigits := parts[1]
		if len(fractionDigits) < 2 {
			fractionDigits = fractionDigits + strings.Repeat("0", 2-len(fractionDigits))
		} else {
			fractionDigits = fractionDigits[:2]
		}
		fraction, err = strconv.ParseInt(fractionDigits, 10, 64)
		if err != nil {
			return domain.Amount{}, false, fmt.Errorf("parse fraction part %q: %w", raw, err)
		}
	}
	cents := whole*100 + fraction
	if negative {
		cents = -cents
	}
	return domain.AmountFromCents(cents), true, nil
}