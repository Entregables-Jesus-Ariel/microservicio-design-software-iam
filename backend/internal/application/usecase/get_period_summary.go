package usecase

import (
	"context"
	"time"

	"ferreteria/internal/application/port"
	"ferreteria/internal/domain"
)

// GetPeriodSummary computes totals for a day, week or month.
type GetPeriodSummary struct {
	movements port.MovementRepository
}

// NewGetPeriodSummary wires the use case with its ports.
func NewGetPeriodSummary(movements port.MovementRepository) *GetPeriodSummary {
	return &GetPeriodSummary{movements: movements}
}

// ExecuteForGranularity summarises the period containing reference.
func (u *GetPeriodSummary) ExecuteForGranularity(
	ctx context.Context,
	granularity domain.PeriodGranularity,
	reference time.Time,
) (domain.PeriodBalance, error) {
	period, err := domain.PeriodFor(granularity, reference)
	if err != nil {
		return domain.PeriodBalance{}, err
	}
	return u.execute(ctx, period)
}

// ExecuteForRange summarises an explicit date range.
func (u *GetPeriodSummary) ExecuteForRange(
	ctx context.Context,
	start time.Time,
	end time.Time,
) (domain.PeriodBalance, error) {
	period, err := domain.NewPeriod(start, end)
	if err != nil {
		return domain.PeriodBalance{}, err
	}
	return u.execute(ctx, period)
}

// execute aggregates totals. Cancelled movements are excluded by the
// repository query, so they never reach the balance.
func (u *GetPeriodSummary) execute(ctx context.Context, period domain.Period) (domain.PeriodBalance, error) {
	totals, err := u.movements.TotalsByCategoryType(ctx, period)
	if err != nil {
		return domain.PeriodBalance{}, err
	}

	balance := domain.PeriodBalance{Period: period}
	for _, total := range totals {
		switch total.Type {
		case domain.CategoryIncome:
			balance.TotalIncome = domain.AmountFromCents(total.Cents)
		case domain.CategoryExpense:
			balance.TotalExpense = domain.AmountFromCents(total.Cents)
		}
	}
	return balance, nil
}
