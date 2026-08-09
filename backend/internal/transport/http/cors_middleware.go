package http

import (
	"net/http"

	"ferreteria/internal/config"
)

// withCORS allows exactly the configured origin. A wildcard would let any
// site call the API with the admin's browser credentials.
func withCORS(settings config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		origin := request.Header.Get("Origin")
		if origin != "" && origin == settings.CORSOrigin {
			header := writer.Header()
			header.Set("Access-Control-Allow-Origin", settings.CORSOrigin)
			header.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			header.Set("Vary", "Origin")
		}

		if request.Method == http.MethodOptions {
			if origin != settings.CORSOrigin {
				writer.WriteHeader(http.StatusForbidden)
				return
			}
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
