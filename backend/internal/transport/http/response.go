package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"iam/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// writeDomainError maps known domain errors to HTTP status codes and
// falls back to 500 for anything unexpected (never leaking internals).
func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailAlreadyRegistered):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrInvalidToken):
		writeError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrAccountLocked):
		writeError(w, http.StatusForbidden, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
