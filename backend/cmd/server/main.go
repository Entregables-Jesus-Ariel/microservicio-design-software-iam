// Command server boots the backend. It only wires dependencies and hands
// control to the transport layer; no business rule or SQL lives here.
package main

import (
	"log"

	"iam/internal/config"
	transporthttp "iam/internal/transport/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	log.Printf("IAM backend listening on port %d", cfg.HTTPPort)
	if err := transporthttp.Run(cfg); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
