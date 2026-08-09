// Package security implements the cryptographic ports declared by the
// application layer. No other layer imports a crypto library directly.
package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"iam/internal/application/port"
)

// JWTTokenService issues signed access tokens and opaque refresh tokens.
type JWTTokenService struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewJWTTokenService builds the service with the signing secret and TTLs.
func NewJWTTokenService(secret string, accessTokenTTL time.Duration) *JWTTokenService {
	return &JWTTokenService{
		secret:          []byte(secret),
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: 7 * 24 * time.Hour,
	}
}

// GenerateAccessToken signs a short-lived JWT carrying the user identity.
func (s *JWTTokenService) GenerateAccessToken(claims port.TokenClaims) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   claims.UserID,
		"email": claims.Email,
		"iat":   now.Unix(),
		"exp":   now.Add(s.accessTokenTTL).Unix(),
	})
	return token.SignedString(s.secret)
}

// GenerateRefreshToken returns a random opaque token plus its SHA-256 hash.
// Only the hash is ever persisted; the plain value is handed to the client once.
func (s *JWTTokenService) GenerateRefreshToken() (string, string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", err
	}
	plain := hex.EncodeToString(raw)

	sum := sha256.Sum256([]byte(plain))
	hash := hex.EncodeToString(sum[:])

	return plain, hash, nil
}
