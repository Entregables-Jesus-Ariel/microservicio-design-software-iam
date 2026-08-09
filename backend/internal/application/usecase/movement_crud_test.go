package usecase

import (
	"context"
	"testing"
	"time"

	"ferreteria/internal/application/port"
	"ferreteria/internal/domain"
)

// TestMovementCrudCycle walks create, read, update and cancel in one flow.
// Declaring CRUD coverage without exercising every operation would let a
// partial implementation reach the frontend.
func TestMovementCrudCycle(t *testing.T) {
	ctx := context.Background()
	movements := newFakeMovementRepository()
	categories := newFakeCategoryRepository()
	audits := &fakeAuditRepository{}
	income := seedCategory(t, categories, "Sales", domain.CategoryIncome)
	expense := seedCategory(t, categories, "Transport", domain.CategoryExpense)
	date := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC)

	created, err := NewRecordMovement(movements, categories, audits).Execute(ctx, RecordMovementInput{
		UserID:      1,
		CategoryID:  income.ID,
		Date:        date,
		AmountCents: 25000,
		Note:        "daily sales",
	})
	if err != nil {
		t.Fatalf("create returned error: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("a created movement must receive an identifier")
	}

	listed, err := NewListMovement(movements).Execute(ctx, ListMovementInput{
		Start: date.AddDate(0, 0, -1),
		End:   date.AddDate(0, 0, 1),
	})
	if err != nil {
		t.Fatalf("list returned error: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("expected one movement in the period, got %d", len(listed))
	}

	edited, err := NewEditMovement(movements, categories, audits).Execute(ctx, EditMovementInput{
		MovementID:  created.ID,
		EditorID:    1,
		CategoryID:  expense.ID,
		Date:        date,
		AmountCents: 40000,
		Note:        "supplier payment",
	})
	if err != nil {
		t.Fatalf("edit returned error: %v", err)
	}
	if edited.Amount.Cents() != 40000 || edited.CategoryID != expense.ID {
		t.Fatalf("edit did not apply: %d cents, category %d", edited.Amount.Cents(), edited.CategoryID)
	}

	if err := NewCancelMovement(movements, audits).Execute(ctx, created.ID, 1); err != nil {
		t.Fatalf("cancel returned error: %v", err)
	}
	stored, err := movements.FindByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("find after cancel returned error: %v", err)
	}
	if !stored.IsCancelled() {
		t.Fatal("a cancelled movement must remain stored and marked")
	}

	if len(audits.entries) != 3 {
		t.Fatalf("expected create, update and delete audit entries, got %d", len(audits.entries))
	}
}

// TestCategoryCrudCycle covers the category operations the API exposes.
func TestCategoryCrudCycle(t *testing.T) {
	ctx := context.Background()
	categories := newFakeCategoryRepository()

	created, err := NewCreateCategory(categories).Execute(ctx, "Services", domain.CategoryExpense)
	if err != nil {
		t.Fatalf("create category returned error: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("a created category must receive an identifier")
	}

	all, err := NewListCategory(categories).Execute(ctx, "")
	if err != nil {
		t.Fatalf("list categories returned error: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected one category, got %d", len(all))
	}

	if _, err := NewCreateCategory(categories).Execute(ctx, "Services", domain.CategoryExpense); err == nil {
		t.Fatal("a duplicate category name must be rejected")
	}
}

// TestSummaryExcludesCancelledMovements pins the reporting rule that a
// cancelled entry must not affect the balance.
func TestSummaryExcludesCancelledMovements(t *testing.T) {
	ctx := context.Background()
	movements := newFakeMovementRepository()
	movements.totals = []port.CategoryTotal{
		{Type: domain.CategoryIncome, Cents: 25000},
	}

	balance, err := NewGetPeriodSummary(movements).ExecuteForGranularity(
		ctx,
		domain.PeriodMonthly,
		time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("summary returned error: %v", err)
	}
	if balance.Net() != 25000 {
		t.Fatalf("expected net 25000, got %d", balance.Net())
	}
}
