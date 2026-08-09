package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewPeriodRejectsInvertedRange(t *testing.T) {
	start := time.Date(2026, time.April, 20, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, time.April, 10, 0, 0, 0, 0, time.UTC)

	if _, err := NewPeriod(start, end); !errors.Is(err, ErrInvalidPeriod) {
		t.Fatalf("expected ErrInvalidPeriod, got %v", err)
	}
}

func TestPeriodForMonthlyCoversWholeMonth(t *testing.T) {
	reference := time.Date(2026, time.April, 15, 13, 30, 0, 0, time.UTC)

	period, err := PeriodFor(PeriodMonthly, reference)
	if err != nil {
		t.Fatalf("PeriodFor returned error: %v", err)
	}
	if period.Start.Day() != 1 || period.Start.Month() != time.April {
		t.Fatalf("unexpected period start: %v", period.Start)
	}
	if period.End.Day() != 30 {
		t.Fatalf("expected April to end on day 30, got %d", period.End.Day())
	}
	if !period.Contains(reference) {
		t.Fatal("the reference date must fall inside its own month")
	}
}

func TestPeriodForWeeklyStartsOnMonday(t *testing.T) {
	// 2026-04-15 is a Wednesday.
	reference := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC)

	period, err := PeriodFor(PeriodWeekly, reference)
	if err != nil {
		t.Fatalf("PeriodFor returned error: %v", err)
	}
	if period.Start.Weekday() != time.Monday {
		t.Fatalf("expected week to start on Monday, got %v", period.Start.Weekday())
	}
	if period.End.Weekday() != time.Sunday {
		t.Fatalf("expected week to end on Sunday, got %v", period.End.Weekday())
	}
}

func TestPeriodForRejectsUnknownGranularity(t *testing.T) {
	reference := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC)

	if _, err := PeriodFor(PeriodGranularity("yearly"), reference); !errors.Is(err, ErrInvalidPeriod) {
		t.Fatalf("expected ErrInvalidPeriod, got %v", err)
	}
}

func TestPeriodBalanceNetSubtractsExpense(t *testing.T) {
	period, err := PeriodFor(PeriodDaily, time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("PeriodFor returned error: %v", err)
	}
	balance := PeriodBalance{
		Period:       period,
		TotalIncome:  AmountFromCents(150000),
		TotalExpense: AmountFromCents(40000),
	}

	if balance.Net() != 110000 {
		t.Fatalf("expected net 110000, got %d", balance.Net())
	}
}
