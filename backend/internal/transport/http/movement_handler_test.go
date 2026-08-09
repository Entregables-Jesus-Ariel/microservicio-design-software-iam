package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateMovementRejectsRequestWithoutToken(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/movements",
		strings.NewReader(`{"category_id":1,"date":"2026-04-15","amount_cents":2500}`),
	)

	handler := withAuth(newTestTokenService(t), http.HandlerFunc(
		func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusCreated)
		},
	))
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without a token, got %d", recorder.Code)
	}
}

func TestCreateMovementRejectsMalformedBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/movements", strings.NewReader("{not json"))

	handleCreateMovement(&dependencies{})(recorder, injectCaller(request))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for malformed JSON, got %d", recorder.Code)
	}
}

func TestCreateMovementRejectsNonIsoDate(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/movements",
		strings.NewReader(`{"category_id":1,"date":"15-04-2026","amount_cents":2500}`),
	)

	handleCreateMovement(&dependencies{})(recorder, injectCaller(request))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for a non-ISO date, got %d", recorder.Code)
	}
}

func TestMovementIDFromPathReadsTrailingIdentifier(t *testing.T) {
	id, err := movementIDFromPath("/api/movements/42")
	if err != nil {
		t.Fatalf("movementIDFromPath returned error: %v", err)
	}
	if id != 42 {
		t.Fatalf("expected 42, got %d", id)
	}
}
