package http

import (
	"encoding/json"
	"net/http"
)

// handleLogin authenticates the admin and returns an access token.
func handleLogin(dependencies *dependencies) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var payload loginRequest
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			writeError(writer, http.StatusBadRequest, "request body is not valid JSON")
			return
		}

		token, err := dependencies.authenticateUser.Execute(
			request.Context(),
			payload.Username,
			payload.Password,
		)
		if err != nil {
			writeDomainError(writer, request, err)
			return
		}
		writeJSON(writer, http.StatusOK, tokenResponse{Token: token})
	}
}
