package http

import "net/http"

func registerAuthRoutes(mux *http.ServeMux, deps *dependencies) {
	mux.HandleFunc("POST /api/auth/register", handleRegister(deps))
	mux.HandleFunc("POST /api/auth/login", handleLogin(deps))
	mux.HandleFunc("POST /api/auth/refresh", handleRefresh(deps))
	mux.HandleFunc("POST /api/auth/logout", handleLogout(deps))
	mux.HandleFunc("POST /api/auth/forgot-password", handleForgotPassword(deps))
	mux.HandleFunc("POST /api/auth/reset-password", handleResetPassword(deps))
}
