package http

import (
	"encoding/json"
	"net/http"

	"ferreteria/internal/domain"
)

// handleListCategory returns categories, optionally filtered by type.
func handleListCategory(dependencies *dependencies) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		filter := domain.CategoryType(request.URL.Query().Get("type"))
		categories, err := dependencies.listCategory.Execute(request.Context(), filter)
		if err != nil {
			writeDomainError(writer, request, err)
			return
		}
		writeJSON(writer, http.StatusOK, newCategoryResponses(categories))
	}
}

// handleCreateCategory adds a custom category.
func handleCreateCategory(dependencies *dependencies) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var payload categoryRequest
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			writeError(writer, http.StatusBadRequest, "request body is not valid JSON")
			return
		}

		category, err := dependencies.createCategory.Execute(
			request.Context(),
			payload.Name,
			domain.CategoryType(payload.Type),
		)
		if err != nil {
			writeDomainError(writer, request, err)
			return
		}
		writeJSON(writer, http.StatusCreated, newCategoryResponse(category))
	}
}
