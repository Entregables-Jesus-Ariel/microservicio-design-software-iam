package http

import "net/http"

func registerRBACRoutes(mux *http.ServeMux, deps *dependencies, secret string) {
	// Protected by auth AND must be ADMIN
	adminOnly := func(h http.HandlerFunc) http.HandlerFunc {
		return RequireAuth(secret, RequireRole("ADMIN", h))
	}

	mux.HandleFunc("GET /api/admin/users", adminOnly(handleListUsers(deps)))
	mux.HandleFunc("GET /api/admin/roles", adminOnly(handleListRoles(deps)))
	mux.HandleFunc("POST /api/admin/users/{id}/roles", adminOnly(handleAssignRole(deps)))
	mux.HandleFunc("DELETE /api/admin/users/{id}/roles/{roleId}", adminOnly(handleRevokeRole(deps)))
}
