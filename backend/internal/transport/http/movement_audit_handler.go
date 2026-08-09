package http

import (
	"net/http"
	"strconv"
	"strings"
)

// handleGetMovementAudit returns the change history of a movement.
func handleGetMovementAudit(dependencies *dependencies) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, ok := callerFrom(request.Context())
		if !ok {
			writeError(writer, http.StatusUnauthorized, "authentication required")
			return
		}

		movementID, err := movementIDFromAuditPath(request.URL.Path)
		if err != nil {
			writeError(writer, http.StatusBadRequest, "movement id must be numeric")
			return
		}

		entries, err := dependencies.getMovementAudit.Execute(request.Context(), movementID)
		if err != nil {
			writeDomainError(writer, request, err)
			return
		}
		writeJSON(writer, http.StatusOK, newMovementAuditResponses(entries))
	}
}

// movementIDFromAuditPath extracts the identifier from
// /api/movements/{id}/audit.
func movementIDFromAuditPath(path string) (int64, error) {
	trimmed := strings.TrimSuffix(path, "/")
	trimmed = strings.TrimSuffix(trimmed, "/audit")
	index := strings.LastIndex(trimmed, "/")
	return strconv.ParseInt(trimmed[index+1:], 10, 64)
}