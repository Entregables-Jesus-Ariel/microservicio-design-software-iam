package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"iam/internal/application/usecase"
)

// handleRegister creates a new user account with the default role.
func handleRegister(deps *dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload registerRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "request body is not valid JSON")
			return
		}

		if strings.TrimSpace(payload.Email) == "" || len(payload.Password) < 8 {
			writeError(w, http.StatusBadRequest, "email is required and password must have at least 8 characters")
			return
		}

		user, err := deps.registerUser.Execute(r.Context(), registerInput(payload))
		if err != nil {
			writeDomainError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, userResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		})
	}
}

func registerInput(p registerRequest) usecase.RegisterUserInput {
	return usecase.RegisterUserInput{
		Email:     p.Email,
		Password:  p.Password,
		FirstName: p.FirstName,
		LastName:  p.LastName,
	}
}
