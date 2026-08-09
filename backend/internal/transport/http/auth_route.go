package http

import "net/http"

func registerAuthRoutes(mux *http.ServeMux, deps *dependencies) {
	mux.HandleFunc("POST /api/auth/register", handleRegister(deps))
}
