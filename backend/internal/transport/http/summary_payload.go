package http

import "ferreteria/internal/domain"

// summaryResponse reports the totals of a period. Net may be negative, so
// it is a plain integer rather than a positive-only amount.
type summaryResponse struct {
	Start             string `json:"start"`
	End               string `json:"end"`
	TotalIncomeCents  int64  `json:"total_income_cents"`
	TotalExpenseCents int64  `json:"total_expense_cents"`
	NetCents          int64  `json:"net_cents"`
}

func newSummaryResponse(balance domain.PeriodBalance) summaryResponse {
	return summaryResponse{
		Start:             balance.Period.Start.Format("2006-01-02"),
		End:               balance.Period.End.Format("2006-01-02"),
		TotalIncomeCents:  balance.TotalIncome.Cents(),
		TotalExpenseCents: balance.TotalExpense.Cents(),
		NetCents:          balance.Net(),
	}
}
