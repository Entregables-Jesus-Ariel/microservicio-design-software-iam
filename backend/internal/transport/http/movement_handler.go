package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ferreteria/internal/application/usecase"
)

// handleCreateMovement records an income or expense entry.
func handleCreateMovement(dependencies *dependencies) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		caller, ok := callerFrom(request.Context())
		if !ok {
			writeError(writer, http.StatusUnauthorized, "authentication required")
			return
		}

		var payload movementRequest
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			writeError(writer, http.StatusBadRequest, "request body is not valid JSON")
			return
		}
		date, err := payload.parsedDate()
		if err != nil {
			writeDomainError(writer, request, err)
			return
		}

		movement, err := dependencies.recordMovement.Execute(request.Context(), usecase.RecordMovementInput{
			UserID:      caller.UserID,
			CategoryID:  payload.CategoryID,
			Date:        date,
			AmountCents: payload.AmountCents,
			Note:        payload.Note,
		})
		if err != nil {
			writeDomainError(writer, request, err)
			return
		}
		writeJSON(writer, http.StatusCreated, newMovementResponse(movement))
	}
}

// handleListMovement returns the movements of a period.
func handleListMovement(dependencies *dependencies) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		start, err := time.Parse("2006-01-02", query.Get("start"))
		if err != nil {
			writeError(writer, http.StatusBadRequest, "start must be an ISO date")
			return
		}
		end, err := time.Parse("2006-01-02", query.Get("end"))
		if err != nil {
			writeError(writer, http.StatusBadRequest, "end must be an ISO date")
			return
		}

		categoryID, _ := strconv.ParseInt(query.Get("category_id"), 10, 64)
		page, _ := strconv.Atoi(query.Get("page"))
		pageSize, _ := strconv.Atoi(query.Get("page_size"))

		movements, err := dependencies.listMovement.Execute(request.Context(), usecase.ListMovementInput{
			Start:      start,
			End:        end,
			CategoryID: categoryID,
			Page:       page,
			PageSize:   pageSize,
		})
		if err != nil {
			writeDomainError(writer, request, err)
			return
		}
		writeJSON(writer, http.StatusOK, newMovementResponses(movements))
	}
}

// handleEditMovement updates a movement and records the change.
func handleEditMovement(dependencies *dependencies) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		caller, ok := callerFrom(request.Context())
		if !ok {
			writeError(writer, http.StatusUnauthorized, "authentication required")
			return
		}
		movementID, err := movementIDFromPath(request.URL.Path)
		if err != nil {
			writeError(writer, http.StatusBadRequest, "movement id must be numeric")
			return
		}

		var payload movementRequest
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			writeError(writer, http.StatusBadRequest, "request body is not valid JSON")
			return
		}
		date, err := payload.parsedDate()
		if err != nil {
			writeDomainError(writer, request, err)
			return
		}

		movement, err := dependencies.editMovement.Execute(request.Context(), usecase.EditMovementInput{
			MovementID:  movementID,
			EditorID:    caller.UserID,
			CategoryID:  payload.CategoryID,
			Date:        date,
			AmountCents: payload.AmountCents,
			Note:        payload.Note,
		})
		if err != nil {
			writeDomainError(writer, request, err)
			return
		}
		writeJSON(writer, http.StatusOK, newMovementResponse(movement))
	}
}

// handleCancelMovement soft deletes a movement.
func handleCancelMovement(dependencies *dependencies) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		caller, ok := callerFrom(request.Context())
		if !ok {
			writeError(writer, http.StatusUnauthorized, "authentication required")
			return
		}
		movementID, err := movementIDFromPath(request.URL.Path)
		if err != nil {
			writeError(writer, http.StatusBadRequest, "movement id must be numeric")
			return
		}

		if err := dependencies.cancelMovement.Execute(request.Context(), movementID, caller.UserID); err != nil {
			writeDomainError(writer, request, err)
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}

// movementIDFromPath extracts the trailing identifier of /movements/{id}.
func movementIDFromPath(path string) (int64, error) {
	trimmed := strings.TrimSuffix(path, "/")
	index := strings.LastIndex(trimmed, "/")
	return strconv.ParseInt(trimmed[index+1:], 10, 64)
}
