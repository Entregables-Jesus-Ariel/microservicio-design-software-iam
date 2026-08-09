package http

import (
	"net/http"
	"time"

	"ferreteria/internal/domain"
)

// handleGetSummary returns totals for a granularity or an explicit range.
func handleGetSummary(dependencies *dependencies) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()

		if start := query.Get("start"); start != "" {
			from, err := time.Parse("2006-01-02", start)
			if err != nil {
				writeError(writer, http.StatusBadRequest, "start must be an ISO date")
				return
			}
			to, err := time.Parse("2006-01-02", query.Get("end"))
			if err != nil {
				writeError(writer, http.StatusBadRequest, "end must be an ISO date")
				return
			}
			balance, err := dependencies.getPeriodSummary.ExecuteForRange(request.Context(), from, to)
			if err != nil {
				writeDomainError(writer, request, err)
				return
			}
			writeJSON(writer, http.StatusOK, newSummaryResponse(balance))
			return
		}

		granularity := domain.PeriodGranularity(query.Get("granularity"))
		if granularity == "" {
			granularity = domain.PeriodMonthly
		}
		balance, err := dependencies.getPeriodSummary.ExecuteForGranularity(
			request.Context(),
			granularity,
			time.Now().UTC(),
		)
		if err != nil {
			writeDomainError(writer, request, err)
			return
		}
		writeJSON(writer, http.StatusOK, newSummaryResponse(balance))
	}
}
