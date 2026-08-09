// Package http adapts the application use cases to HTTP. It owns request
// decoding, status mapping and routing, and holds no business rule.
package http

import (
	"fmt"
	"net/http"

	"ferreteria/internal/config"
)

// idleTimeoutFactor derives the idle timeout from the write timeout so a
// slow client cannot hold a connection open indefinitely.
const idleTimeoutFactor = 4

// Run wires dependencies and serves until the process is stopped.
func Run(settings config.Config) error {
	dependencies, err := buildDependencies(settings)
	if err != nil {
		return err
	}
	defer dependencies.Close()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", settings.HTTPPort),
		Handler:      newHandler(settings, dependencies),
		ReadTimeout:  settings.ReadTimeout,
		WriteTimeout: settings.WriteTimeout,
		IdleTimeout:  settings.WriteTimeout * idleTimeoutFactor,
	}
	return server.ListenAndServe()
}

// newHandler assembles the route table and the middleware chain.
func newHandler(settings config.Config, dependencies *dependencies) http.Handler {
	mux := http.NewServeMux()
	registerAuthRoutes(mux, dependencies)
	registerMovementRoutes(mux, settings, dependencies)
	registerCategoryRoutes(mux, settings, dependencies)
	registerSummaryRoutes(mux, settings, dependencies)
	return withCORS(settings, withRecovery(mux))
}
