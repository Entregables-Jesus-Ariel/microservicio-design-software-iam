// Package http adapts the application use cases to HTTP. It owns request
// decoding, status mapping and routing, and holds no business rule.
package http

import (
	"fmt"
	"net/http"

	"iam/internal/config"
)

const idleTimeoutFactor = 4

// Run wires dependencies and serves until the process is stopped.
func Run(cfg config.Config) error {
	deps, err := buildDependencies(cfg)
	if err != nil {
		return err
	}
	defer deps.Close()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:      newHandler(cfg, deps),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.WriteTimeout * idleTimeoutFactor,
	}
	return server.ListenAndServe()
}

func newHandler(cfg config.Config, deps *dependencies) http.Handler {
	mux := http.NewServeMux()
	registerAuthRoutes(mux, deps)
	return withCORS(cfg, withRecovery(mux))
}
