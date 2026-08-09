package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"ferreteria/internal/application/usecase"
	"ferreteria/internal/domain"
)

// errorBody is the only error shape the API returns.
type errorBody struct {
	Error string `json:"error"`
}

// withRecovery turns a panic into a generic 500 so no stack trace, file
// path or SQL fragment ever reaches the client.
func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("unhandled panic on %s %s: %v", request.Method, request.URL.Path, recovered)
				writeError(writer, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(writer, request)
	})
}

// writeError emits a sanitized error payload.
func writeError(writer http.ResponseWriter, status int, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(errorBody{Error: message})
}

// writeJSON emits a successful payload.
func writeJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}

// writeDomainError maps a known failure to its status code. Unknown errors
// are logged server-side and reported generically.
func writeDomainError(writer http.ResponseWriter, request *http.Request, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		writeError(writer, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, domain.ErrMovementNotFound),
		errors.Is(err, domain.ErrCategoryNotFound),
		errors.Is(err, domain.ErrUserNotFound):
		writeError(writer, http.StatusNotFound, "resource not found")
	case errors.Is(err, usecase.ErrCategoryNameTaken):
		writeError(writer, http.StatusConflict, "category name is already in use")
	case errors.Is(err, domain.ErrInvalidAmount),
		errors.Is(err, domain.ErrInvalidDate),
		errors.Is(err, domain.ErrInvalidPeriod),
		errors.Is(err, domain.ErrNoteTooLong),
		errors.Is(err, domain.ErrCategoryRequired),
		errors.Is(err, domain.ErrUserRequired),
		errors.Is(err, domain.ErrInvalidCategoryType),
		errors.Is(err, domain.ErrCategoryNameRequired),
		errors.Is(err, domain.ErrMovementCancelled):
		writeError(writer, http.StatusBadRequest, err.Error())
	default:
		log.Printf("unexpected error on %s %s: %v", request.Method, request.URL.Path, err)
		writeError(writer, http.StatusInternalServerError, "internal server error")
	}
}
