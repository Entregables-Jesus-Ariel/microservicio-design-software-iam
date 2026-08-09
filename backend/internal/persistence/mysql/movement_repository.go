package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ferreteria/internal/application/port"
	"ferreteria/internal/domain"
)

// MovementRepository stores movements in MySQL. Every statement is
// parameterized, so user input never reaches the query text.
type MovementRepository struct {
	database *sql.DB
}

func NewMovementRepository(database *sql.DB) *MovementRepository {
	return &MovementRepository{database: database}
}

func (r *MovementRepository) Create(ctx context.Context, movement *domain.Movement) (*domain.Movement, error) {
	const statement = `
		INSERT INTO movements (user_id, category_id, date, amount, note, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.database.ExecContext(
		ctx,
		statement,
		movement.UserID,
		movement.CategoryID,
		movement.Date,
		centsToDecimalString(movement.Amount.Cents()),
		noteValue(movement.Note),
		movement.CreatedAt,
		movement.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert movement: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("read inserted movement id: %w", err)
	}
	movement.ID = id
	return movement, nil
}

func (r *MovementRepository) Update(ctx context.Context, movement *domain.Movement) error {
	const statement = `
		UPDATE movements
		SET category_id = ?, date = ?, amount = ?,
		    note = ?, updated_at = ?, cancelled_at = ?
		WHERE id = ?`

	result, err := r.database.ExecContext(
		ctx,
		statement,
		movement.CategoryID,
		movement.Date,
		centsToDecimalString(movement.Amount.Cents()),
		noteValue(movement.Note),
		movement.UpdatedAt,
		cancelledValue(movement.CancelledAt),
		movement.ID,
	)
	if err != nil {
		return fmt.Errorf("update movement: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("inspect update result: %w", err)
	}
	if affected == 0 {
		return domain.ErrMovementNotFound
	}
	return nil
}

func (r *MovementRepository) FindByID(ctx context.Context, id int64) (*domain.Movement, error) {
	const query = `
		SELECT id, user_id, category_id, date, CAST(amount * 100 AS SIGNED),
		       note, created_at, updated_at, cancelled_at
		FROM movements
		WHERE id = ?`

	var record movementRecord
	err := r.database.QueryRowContext(ctx, query, id).Scan(
		&record.ID,
		&record.UserID,
		&record.CategoryID,
		&record.Date,
		&record.AmountCents,
		&record.Note,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.CancelledAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrMovementNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("select movement: %w", err)
	}
	return record.toDomain(), nil
}

func (r *MovementRepository) List(ctx context.Context, filter port.MovementFilter) ([]*domain.Movement, error) {
	const query = `
		SELECT id, user_id, category_id, date, CAST(amount * 100 AS SIGNED),
		       note, created_at, updated_at, cancelled_at
		FROM movements
		WHERE date BETWEEN ? AND ?
		  AND (? = 0 OR category_id = ?)
		  AND (? OR cancelled_at IS NULL)
		ORDER BY date DESC, id DESC
		LIMIT ? OFFSET ?`

	rows, err := r.database.QueryContext(
		ctx,
		query,
		filter.Period.Start,
		filter.Period.End,
		filter.CategoryID,
		filter.CategoryID,
		filter.IncludeCancelled,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list movements: %w", err)
	}
	defer rows.Close()

	movements := make([]*domain.Movement, 0)
	for rows.Next() {
		var record movementRecord
		if err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.CategoryID,
			&record.Date,
			&record.AmountCents,
			&record.Note,
			&record.CreatedAt,
			&record.UpdatedAt,
			&record.CancelledAt,
		); err != nil {
			return nil, fmt.Errorf("scan movement row: %w", err)
		}
		movements = append(movements, record.toDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate movement rows: %w", err)
	}
	return movements, nil
}

func (r *MovementRepository) TotalsByCategoryType(
	ctx context.Context,
	period domain.Period,
) ([]port.CategoryTotal, error) {
	const query = `
		SELECT categories.type, CAST(COALESCE(SUM(movements.amount) * 100, 0) AS SIGNED)
		FROM movements
		JOIN categories ON categories.id = movements.category_id
		WHERE movements.date BETWEEN ? AND ?
		  AND movements.cancelled_at IS NULL
		GROUP BY categories.type`

	rows, err := r.database.QueryContext(ctx, query, period.Start, period.End)
	if err != nil {
		return nil, fmt.Errorf("aggregate movement totals: %w", err)
	}
	defer rows.Close()

	totals := make([]port.CategoryTotal, 0, 2)
	for rows.Next() {
		var categoryType string
		var cents int64
		if err := rows.Scan(&categoryType, &cents); err != nil {
			return nil, fmt.Errorf("scan total row: %w", err)
		}
		totals = append(totals, port.CategoryTotal{
			Type:  domain.CategoryType(categoryType),
			Cents: cents,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate total rows: %w", err)
	}
	return totals, nil
}