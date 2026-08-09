package http

import "net/http"

// Authentication endpoint exposed by this router:
//   POST /api/auth/login  authenticate the admin and issue an access token
//
// This is the only route reachable without a bearer token.
func registerAuthRoutes(mux *http.ServeMux, dependencies *dependencies) {
	mux.HandleFunc("/api/auth/login", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		handleLogin(dependencies)(writer, request)
	})
}
