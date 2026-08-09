package domain

import (
	"errors"
	"testing"
	"time"
)

func validAmount(t *testing.T, cents int64) Amount {
	t.Helper()
	amount, err := NewAmount(cents)
	if err != nil {
		t.Fatalf("NewAmount(%d) returned error: %v", cents, err)
	}
	return amount
}

func TestNewMovementRejectsNonPositiveAmount(t *testing.T) {
	if _, err := NewAmount(0); !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
	if _, err := NewAmount(-100); !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount for negative value, got %v", err)
	}
}

func TestNewMovementRejectsMissingRelations(t *testing.T) {
	amount := validAmount(t, 1500)
	date := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC)

	if _, err := NewMovement(0, 2, date, amount, ""); !errors.Is(err, ErrUserRequired) {
		t.Fatalf("expected ErrUserRequired, got %v", err)
	}
	if _, err := NewMovement(1, 0, date, amount, ""); !errors.Is(err, ErrCategoryRequired) {
		t.Fatalf("expected ErrCategoryRequired, got %v", err)
	}
	if _, err := NewMovement(1, 2, time.Time{}, amount, ""); !errors.Is(err, ErrInvalidDate) {
		t.Fatalf("expected ErrInvalidDate, got %v", err)
	}
}

func TestNewMovementRejectsOversizedNote(t *testing.T) {
	amount := validAmount(t, 1500)
	date := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC)
	note := make([]byte, MaxNoteLength+1)
	for index := range note {
		note[index] = 'a'
	}

	if _, err := NewMovement(1, 2, date, amount, string(note)); !errors.Is(err, ErrNoteTooLong) {
		t.Fatalf("expected ErrNoteTooLong, got %v", err)
	}
}

func TestCancelMarksMovementAndRejectsSecondCancel(t *testing.T) {
	amount := validAmount(t, 2500)
	date := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC)
	movement, err := NewMovement(1, 2, date, amount, "daily sales")
	if err != nil {
		t.Fatalf("NewMovement returned error: %v", err)
	}
	if movement.IsCancelled() {
		t.Fatal("a new movement must not be cancelled")
	}

	cancelledAt := time.Date(2026, time.April, 16, 10, 0, 0, 0, time.UTC)
	if err := movement.Cancel(cancelledAt); err != nil {
		t.Fatalf("Cancel returned error: %v", err)
	}
	if !movement.IsCancelled() {
		t.Fatal("movement must be cancelled after Cancel")
	}
	if err := movement.Cancel(cancelledAt); !errors.Is(err, ErrMovementCancelled) {
		t.Fatalf("expected ErrMovementCancelled on second cancel, got %v", err)
	}
}

func TestEditRejectsCancelledMovement(t *testing.T) {
	amount := validAmount(t, 2500)
	date := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC)
	movement, err := NewMovement(1, 2, date, amount, "")
	if err != nil {
		t.Fatalf("NewMovement returned error: %v", err)
	}
	if err := movement.Cancel(date); err != nil {
		t.Fatalf("Cancel returned error: %v", err)
	}

	if err := movement.Edit(3, date, amount, "updated"); !errors.Is(err, ErrMovementCancelled) {
		t.Fatalf("expected ErrMovementCancelled, got %v", err)
	}
}

func TestEditAppliesNewValues(t *testing.T) {
	amount := validAmount(t, 2500)
	date := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC)
	movement, err := NewMovement(1, 2, date, amount, "")
	if err != nil {
		t.Fatalf("NewMovement returned error: %v", err)
	}

	newAmount := validAmount(t, 4000)
	newDate := time.Date(2026, time.April, 20, 0, 0, 0, 0, time.UTC)
	if err := movement.Edit(5, newDate, newAmount, "supplier payment"); err != nil {
		t.Fatalf("Edit returned error: %v", err)
	}
	if movement.CategoryID != 5 {
		t.Fatalf("expected category 5, got %d", movement.CategoryID)
	}
	if movement.Amount.Cents() != 4000 {
		t.Fatalf("expected 4000 cents, got %d", movement.Amount.Cents())
	}
	if movement.Note != "supplier payment" {
		t.Fatalf("unexpected note: %q", movement.Note)
	}
}
