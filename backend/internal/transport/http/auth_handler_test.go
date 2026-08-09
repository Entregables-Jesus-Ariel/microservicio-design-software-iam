package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ferreteria/internal/application/port"
	"ferreteria/internal/security"
)

// newTestTokenService builds a token service with a fixed test secret.
func newTestTokenService(t *testing.T) *security.JWTTokenService {
	t.Helper()
	return security.NewJWTTokenService("test-signing-secret-for-unit-tests")
}

// injectCaller attaches an authenticated caller to a request context.
func injectCaller(request *http.Request) *http.Request {
	claims := port.TokenClaims{
		UserID:    1,
		Username:  "admin",
		Role:      "admin",
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}
	return request.WithContext(context.WithValue(request.Context(), callerKey, claims))
}

func TestLoginRejectsMalformedBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader("{"))

	handleLogin(&dependencies{})(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for malformed JSON, got %d", recorder.Code)
	}
}

func TestIssuedTokenIsAcceptedByTheAuthMiddleware(t *testing.T) {
	tokens := newTestTokenService(t)
	token, err := tokens.Issue(port.TokenClaims{
		UserID:    1,
		Username:  "admin",
		Role:      "admin",
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/summary", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	reached := false
	withAuth(tokens, http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		reached = true
		writer.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, request)

	if !reached {
		t.Fatal("a valid token must reach the protected handler")
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
}

func TestExpiredTokenIsRejected(t *testing.T) {
	tokens := newTestTokenService(t)
	token, err := tokens.Issue(port.TokenClaims{
		UserID:    1,
		Username:  "admin",
		Role:      "admin",
		ExpiresAt: time.Now().UTC().Add(-time.Minute),
	})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/summary", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	withAuth(tokens, http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for an expired token, got %d", recorder.Code)
	}
}
