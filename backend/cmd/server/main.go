// Command server boots the backend. It only wires dependencies and hands
// control to the transport layer; no business rule or SQL lives here.
package main

import (
	"log"

	"ferreteria/internal/config"
	transporthttp "ferreteria/internal/transport/http"
)

func main() {
	settings, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	if err := transporthttp.Run(settings); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
