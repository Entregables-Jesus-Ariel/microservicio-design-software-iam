package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ferreteria/internal/domain"
)

func TestPanicBecomesGenericInternalError(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/summary", nil)

	withRecovery(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("pq: relation \"movement\" does not exist")
	})).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", recorder.Code)
	}
	body := recorder.Body.String()
	for _, leak := range []string{"pq:", "relation", "panic", "goroutine"} {
		if strings.Contains(body, leak) {
			t.Fatalf("error response leaked %q: %s", leak, body)
		}
	}
}

func TestUnknownErrorIsReportedGenerically(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/summary", nil)

	writeDomainError(recorder, request, errors.New("dial tcp 127.0.0.1:5432: connection refused"))

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", recorder.Code)
	}
	var payload errorBody
	if err := json.NewDecoder(recorder.Body).Decode(&payload); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if payload.Error != "internal server error" {
		t.Fatalf("unexpected error message: %q", payload.Error)
	}
}

func TestInvalidCredentialsMapToUnauthorized(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)

	writeDomainError(recorder, request, domain.ErrInvalidCredentials)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestNotFoundMapsToNotFoundStatus(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/movements/9", nil)

	writeDomainError(recorder, request, domain.ErrMovementNotFound)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", recorder.Code)
	}
}
