package domain

import "time"

// MaxNoteLength bounds the optional note so a single record cannot carry
// an unbounded payload of potentially personal data.
const MaxNoteLength = 500

// Amount is a monetary value stored as whole cents to avoid the rounding
// drift a floating point type would introduce in financial totals.
type Amount struct {
	cents int64
}

// NewAmount builds a positive monetary amount from cents.
func NewAmount(cents int64) (Amount, error) {
	if cents <= 0 {
		return Amount{}, ErrInvalidAmount
	}
	return Amount{cents: cents}, nil
}

// AmountFromCents rebuilds an amount already validated by persistence.
func AmountFromCents(cents int64) Amount {
	return Amount{cents: cents}
}

// Cents exposes the raw value for persistence and calculations.
func (a Amount) Cents() int64 {
	return a.cents
}

// Add returns the sum of two amounts.
func (a Amount) Add(other Amount) Amount {
	return Amount{cents: a.cents + other.cents}
}

// Movement is a single income or expense entry.
type Movement struct {
	ID          int64
	UserID      int64
	CategoryID  int64
	Date        time.Time
	Amount      Amount
	Note        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CancelledAt *time.Time
}

// NewMovement validates and builds a movement entry.
func NewMovement(userID int64, categoryID int64, date time.Time, amount Amount, note string) (*Movement, error) {
	if userID <= 0 {
		return nil, ErrUserRequired
	}
	if categoryID <= 0 {
		return nil, ErrCategoryRequired
	}
	if date.IsZero() {
		return nil, ErrInvalidDate
	}
	if amount.Cents() <= 0 {
		return nil, ErrInvalidAmount
	}
	if len(note) > MaxNoteLength {
		return nil, ErrNoteTooLong
	}
	now := time.Now().UTC()
	return &Movement{
		UserID:     userID,
		CategoryID: categoryID,
		Date:       date,
		Amount:     amount,
		Note:       note,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// IsCancelled reports whether the movement was soft deleted.
func (m *Movement) IsCancelled() bool {
	return m.CancelledAt != nil
}

// Cancel marks the movement as cancelled without discarding the record.
func (m *Movement) Cancel(at time.Time) error {
	if m.IsCancelled() {
		return ErrMovementCancelled
	}
	cancelled := at.UTC()
	m.CancelledAt = &cancelled
	m.UpdatedAt = cancelled
	return nil
}

// Edit applies new values after validating them.
func (m *Movement) Edit(categoryID int64, date time.Time, amount Amount, note string) error {
	if m.IsCancelled() {
		return ErrMovementCancelled
	}
	if categoryID <= 0 {
		return ErrCategoryRequired
	}
	if date.IsZero() {
		return ErrInvalidDate
	}
	if amount.Cents() <= 0 {
		return ErrInvalidAmount
	}
	if len(note) > MaxNoteLength {
		return ErrNoteTooLong
	}
	m.CategoryID = categoryID
	m.Date = date
	m.Amount = amount
	m.Note = note
	m.UpdatedAt = time.Now().UTC()
	return nil
}
