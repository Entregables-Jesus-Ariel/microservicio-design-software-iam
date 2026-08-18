package http

import (
	"net/http"
	"strconv"
)

func handleListAuditLogins(deps *dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		page, _ := strconv.Atoi(pageStr)
		limit, _ := strconv.Atoi(limitStr)

		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 20
		}

		paginatedLogs, err := deps.auditUsecases.ListAuditLogins(r.Context(), page, limit)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to list audit logs")
			return
		}

		writeJSON(w, http.StatusOK, paginatedLogs)
	}
}
