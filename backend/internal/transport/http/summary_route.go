package http

import (
	"net/http"

	"ferreteria/internal/config"
)

// Summary endpoint exposed by this router:
//   GET /api/summary  totals for a granularity or an explicit date range
func registerSummaryRoutes(
	mux *http.ServeMux,
	settings config.Config,
	dependencies *dependencies,
) {
	summary := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		handleGetSummary(dependencies)(writer, request)
	})

	_ = settings
	mux.Handle("/api/summary", withAuth(dependencies.tokens, summary))
}
