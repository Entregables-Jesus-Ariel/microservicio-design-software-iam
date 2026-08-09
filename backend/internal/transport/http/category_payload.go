package http

import "ferreteria/internal/domain"

// categoryRequest is the accepted shape for creating a category.
type categoryRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// categoryResponse is the category shape returned to clients.
type categoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func newCategoryResponse(category *domain.Category) categoryResponse {
	return categoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Type: string(category.Type),
	}
}

func newCategoryResponses(categories []*domain.Category) []categoryResponse {
	responses := make([]categoryResponse, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, newCategoryResponse(category))
	}
	return responses
}
