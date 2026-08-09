package domain

import "time"

// PeriodGranularity names the supported summary windows.
type PeriodGranularity string

// Supported granularities for period summaries.
const (
	PeriodDaily   PeriodGranularity = "daily"
	PeriodWeekly  PeriodGranularity = "weekly"
	PeriodMonthly PeriodGranularity = "monthly"
)

// Period is an inclusive date range used to summarise movements.
type Period struct {
	Start time.Time
	End   time.Time
}

// NewPeriod validates an explicit date range.
func NewPeriod(start time.Time, end time.Time) (Period, error) {
	if start.IsZero() || end.IsZero() {
		return Period{}, ErrInvalidDate
	}
	if end.Before(start) {
		return Period{}, ErrInvalidPeriod
	}
	return Period{Start: startOfDay(start), End: endOfDay(end)}, nil
}

// PeriodFor builds the range covering reference for a granularity.
func PeriodFor(granularity PeriodGranularity, reference time.Time) (Period, error) {
	if reference.IsZero() {
		return Period{}, ErrInvalidDate
	}
	switch granularity {
	case PeriodDaily:
		return Period{Start: startOfDay(reference), End: endOfDay(reference)}, nil
	case PeriodWeekly:
		weekday := (int(reference.Weekday()) + 6) % 7
		start := reference.AddDate(0, 0, -weekday)
		return Period{Start: startOfDay(start), End: endOfDay(start.AddDate(0, 0, 6))}, nil
	case PeriodMonthly:
		start := time.Date(reference.Year(), reference.Month(), 1, 0, 0, 0, 0, time.UTC)
		return Period{Start: start, End: endOfDay(start.AddDate(0, 1, -1))}, nil
	default:
		return Period{}, ErrInvalidPeriod
	}
}

// Contains reports whether a date falls inside the period.
func (p Period) Contains(date time.Time) bool {
	return !date.Before(p.Start) && !date.After(p.End)
}

// PeriodBalance holds the totals for one period.
type PeriodBalance struct {
	Period       Period
	TotalIncome  Amount
	TotalExpense Amount
}

// Net returns income minus expense in cents. It may be negative, so it is
// not expressed as an Amount, which only accepts positive values.
func (b PeriodBalance) Net() int64 {
	return b.TotalIncome.Cents() - b.TotalExpense.Cents()
}

func startOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}

func endOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 23, 59, 59, 0, time.UTC)
}
