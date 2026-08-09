package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"ferreteria/internal/domain"
)

func TestCancelMovementMarksRecordAndAppendsAudit(t *testing.T) {
	movements := newFakeMovementRepository()
	categories := newFakeCategoryRepository()
	audits := &fakeAuditRepository{}
	category := seedCategory(t, categories, "Sales", domain.CategoryIncome)

	recorded, err := NewRecordMovement(movements, categories, audits).Execute(
		context.Background(),
		RecordMovementInput{
			UserID:      1,
			CategoryID:  category.ID,
			Date:        time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC),
			AmountCents: 25000,
		},
	)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if err := NewCancelMovement(movements, audits).Execute(context.Background(), recorded.ID, 1); err != nil {
		t.Fatalf("Cancel returned error: %v", err)
	}

	stored, err := movements.FindByID(context.Background(), recorded.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if !stored.IsCancelled() {
		t.Fatal("movement must be cancelled")
	}
	if len(audits.entries) != 2 {
		t.Fatalf("expected create and delete audit entries, got %d", len(audits.entries))
	}
	if audits.entries[1].Action != domain.AuditDelete {
		t.Fatalf("expected delete action, got %s", audits.entries[1].Action)
	}
}

func TestCancelMovementRejectsUnknownMovement(t *testing.T) {
	movements := newFakeMovementRepository()
	audits := &fakeAuditRepository{}

	err := NewCancelMovement(movements, audits).Execute(context.Background(), 404, 1)
	if !errors.Is(err, domain.ErrMovementNotFound) {
		t.Fatalf("expected ErrMovementNotFound, got %v", err)
	}
}

func TestCancelMovementRejectsSecondCancellation(t *testing.T) {
	movements := newFakeMovementRepository()
	categories := newFakeCategoryRepository()
	audits := &fakeAuditRepository{}
	category := seedCategory(t, categories, "Sales", domain.CategoryIncome)

	recorded, err := NewRecordMovement(movements, categories, audits).Execute(
		context.Background(),
		RecordMovementInput{
			UserID:      1,
			CategoryID:  category.ID,
			Date:        time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC),
			AmountCents: 25000,
		},
	)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	useCase := NewCancelMovement(movements, audits)
	if err := useCase.Execute(context.Background(), recorded.ID, 1); err != nil {
		t.Fatalf("first cancel returned error: %v", err)
	}
	if err := useCase.Execute(context.Background(), recorded.ID, 1); !errors.Is(err, domain.ErrMovementCancelled) {
		t.Fatalf("expected ErrMovementCancelled, got %v", err)
	}
}
