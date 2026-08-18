package http

import "net/http"

func registerAuditRoutes(mux *http.ServeMux, deps *dependencies, secret string) {
	// Protected by auth AND must be ADMIN
	adminOnly := func(h http.HandlerFunc) http.HandlerFunc {
		return RequireAuth(secret, RequireRole("ADMIN", h))
	}

	mux.HandleFunc("GET /api/admin/audit/logins", adminOnly(handleListAuditLogins(deps)))
}
