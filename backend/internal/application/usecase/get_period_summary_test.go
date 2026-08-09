package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"ferreteria/internal/application/port"
	"ferreteria/internal/domain"
)

func TestGetPeriodSummaryComputesNetBalance(t *testing.T) {
	movements := newFakeMovementRepository()
	movements.totals = []port.CategoryTotal{
		{Type: domain.CategoryIncome, Cents: 150000},
		{Type: domain.CategoryExpense, Cents: 40000},
	}

	balance, err := NewGetPeriodSummary(movements).ExecuteForGranularity(
		context.Background(),
		domain.PeriodMonthly,
		time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("ExecuteForGranularity returned error: %v", err)
	}
	if balance.TotalIncome.Cents() != 150000 {
		t.Fatalf("expected income 150000, got %d", balance.TotalIncome.Cents())
	}
	if balance.TotalExpense.Cents() != 40000 {
		t.Fatalf("expected expense 40000, got %d", balance.TotalExpense.Cents())
	}
	if balance.Net() != 110000 {
		t.Fatalf("expected net 110000, got %d", balance.Net())
	}
}

func TestGetPeriodSummaryRejectsInvertedRange(t *testing.T) {
	movements := newFakeMovementRepository()

	_, err := NewGetPeriodSummary(movements).ExecuteForRange(
		context.Background(),
		time.Date(2026, time.April, 20, 0, 0, 0, 0, time.UTC),
		time.Date(2026, time.April, 10, 0, 0, 0, 0, time.UTC),
	)
	if !errors.Is(err, domain.ErrInvalidPeriod) {
		t.Fatalf("expected ErrInvalidPeriod, got %v", err)
	}
}

func TestGetPeriodSummaryReturnsZeroWhenNoMovements(t *testing.T) {
	movements := newFakeMovementRepository()

	balance, err := NewGetPeriodSummary(movements).ExecuteForGranularity(
		context.Background(),
		domain.PeriodDaily,
		time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("ExecuteForGranularity returned error: %v", err)
	}
	if balance.Net() != 0 {
		t.Fatalf("expected zero net balance, got %d", balance.Net())
	}
}
