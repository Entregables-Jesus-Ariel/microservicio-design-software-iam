package port

import "time"

// TokenClaims carries the identity encoded in an access token.
type TokenClaims struct {
	UserID    int64
	Username  string
	Role      string
	ExpiresAt time.Time
}

// TokenService issues and validates access tokens.
type TokenService interface {
	Issue(claims TokenClaims) (string, error)
	Validate(token string) (TokenClaims, error)
}
