package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"iam/internal/application/port"
)

type contextKey string

const (
	userClaimsKey contextKey = "userClaims"
)

// RequireAuth ensures a valid JWT is present in the Authorization header.
func RequireAuth(secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "missing or invalid authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		userID, _ := claims["sub"].(string)
		email, _ := claims["email"].(string)

		var roles []string
		if rolesInterface, ok := claims["roles"].([]interface{}); ok {
			for _, r := range rolesInterface {
				if roleStr, ok := r.(string); ok {
					roles = append(roles, roleStr)
				}
			}
		}

		tokenClaims := port.TokenClaims{
			UserID: userID,
			Email:  email,
			Roles:  roles,
		}

		ctx := context.WithValue(r.Context(), userClaimsKey, tokenClaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// RequireRole ensures the authenticated user has a specific role.
// It must be used after RequireAuth.
func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(userClaimsKey).(port.TokenClaims)
		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		hasRole := false
		for _, userRole := range claims.Roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			writeError(w, http.StatusForbidden, "forbidden: insufficient permissions")
			return
		}

		next.ServeHTTP(w, r)
	}
}
