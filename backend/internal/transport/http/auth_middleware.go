package http

import (
	"context"
	"net/http"
	"strings"

	"ferreteria/internal/application/port"
	"ferreteria/internal/security"
)

type contextKey string

// callerKey carries the authenticated caller through the request context.
const callerKey contextKey = "caller"

// withAuth rejects a request whose bearer token is absent or invalid.
func withAuth(tokens *security.JWTTokenService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		header := request.Header.Get("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer"))
		if header == "" || token == "" {
			writeError(writer, http.StatusUnauthorized, "authentication required")
			return
		}

		claims, err := tokens.Validate(token)
		if err != nil {
			writeError(writer, http.StatusUnauthorized, "authentication required")
			return
		}

		ctx := context.WithValue(request.Context(), callerKey, claims)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// callerFrom returns the authenticated caller stored by withAuth.
func callerFrom(ctx context.Context) (port.TokenClaims, bool) {
	claims, ok := ctx.Value(callerKey).(port.TokenClaims)
	return claims, ok
}
