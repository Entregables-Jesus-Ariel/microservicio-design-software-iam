package domain

import (
	"errors"
	"testing"
)

func TestNewCategoryRejectsEmptyName(t *testing.T) {
	if _, err := NewCategory("   ", CategoryIncome); !errors.Is(err, ErrCategoryNameRequired) {
		t.Fatalf("expected ErrCategoryNameRequired, got %v", err)
	}
}

func TestNewCategoryRejectsUnknownType(t *testing.T) {
	if _, err := NewCategory("Sales", CategoryType("transfer")); !errors.Is(err, ErrInvalidCategoryType) {
		t.Fatalf("expected ErrInvalidCategoryType, got %v", err)
	}
}

func TestNewCategoryTrimsNameAndClassifies(t *testing.T) {
	category, err := NewCategory("  Sales  ", CategoryIncome)
	if err != nil {
		t.Fatalf("NewCategory returned error: %v", err)
	}
	if category.Name != "Sales" {
		t.Fatalf("expected trimmed name, got %q", category.Name)
	}
	if !category.IsIncome() {
		t.Fatal("an income category must report IsIncome")
	}

	expense, err := NewCategory("Transport", CategoryExpense)
	if err != nil {
		t.Fatalf("NewCategory returned error: %v", err)
	}
	if expense.IsIncome() {
		t.Fatal("an expense category must not report IsIncome")
	}
}
