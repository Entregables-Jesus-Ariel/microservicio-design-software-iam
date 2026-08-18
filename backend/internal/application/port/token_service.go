package port

import "time"

// TokenClaims is what gets embedded inside the access token.
type TokenClaims struct {
	UserID string
	Email  string
	Roles  []string
}

// GeneratedTokens is the pair returned after a successful login.
type GeneratedTokens struct {
	AccessToken      string
	RefreshToken     string
	RefreshTokenHash string
	RefreshExpiresAt time.Time
}

// TokenService issues and validates JWT access tokens and opaque refresh tokens.
type TokenService interface {
	GenerateAccessToken(claims TokenClaims) (string, error)
	GenerateRefreshToken() (plain string, hash string, err error)
	HashToken(plain string) string
}
