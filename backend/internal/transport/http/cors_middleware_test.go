package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ferreteria/internal/config"
)

func corsSettings() config.Config {
	return config.Config{CORSOrigin: "http://localhost:4200"}
}

func TestAllowedOriginReceivesCORSHeader(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/summary", nil)
	request.Header.Set("Origin", "http://localhost:4200")

	withCORS(corsSettings(), http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, request)

	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:4200" {
		t.Fatalf("expected the configured origin, got %q", got)
	}
}

func TestUnknownOriginReceivesNoCORSHeader(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/summary", nil)
	request.Header.Set("Origin", "http://evil.example")

	withCORS(corsSettings(), http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, request)

	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("an unknown origin must receive no CORS header, got %q", got)
	}
}

func TestPreflightFromUnknownOriginIsForbidden(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/api/summary", nil)
	request.Header.Set("Origin", "http://evil.example")

	withCORS(corsSettings(), http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for a preflight from an unknown origin, got %d", recorder.Code)
	}
}
