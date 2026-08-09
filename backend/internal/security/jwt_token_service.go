package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"ferreteria/internal/application/port"
)

// Token validation errors.
var (
	ErrMalformedToken = errors.New("token is malformed")
	ErrInvalidSignature = errors.New("token signature is invalid")
	ErrExpiredToken   = errors.New("token has expired")
)

// JWTTokenService issues and validates HS256 tokens using only the
// standard library, so the backend carries no JWT dependency.
type JWTTokenService struct {
	secret []byte
}

// NewJWTTokenService builds the service with a signing secret.
func NewJWTTokenService(secret string) *JWTTokenService {
	return &JWTTokenService{secret: []byte(secret)}
}

type tokenPayload struct {
	UserID    int64  `json:"sub"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"exp"`
}

// Issue returns a signed token carrying the supplied claims.
func (s *JWTTokenService) Issue(claims port.TokenClaims) (string, error) {
	header, err := encodeSegment(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}
	payload, err := encodeSegment(tokenPayload{
		UserID:    claims.UserID,
		Username:  claims.Username,
		Role:      claims.Role,
		ExpiresAt: claims.ExpiresAt.UTC().Unix(),
	})
	if err != nil {
		return "", err
	}

	signingInput := header + "." + payload
	return signingInput + "." + s.sign(signingInput), nil
}

// Validate verifies the signature and expiry, then returns the claims.
func (s *JWTTokenService) Validate(token string) (port.TokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return port.TokenClaims{}, ErrMalformedToken
	}

	signingInput := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(s.sign(signingInput)), []byte(parts[2])) {
		return port.TokenClaims{}, ErrInvalidSignature
	}

	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return port.TokenClaims{}, ErrMalformedToken
	}
	var payload tokenPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return port.TokenClaims{}, ErrMalformedToken
	}

	expiresAt := time.Unix(payload.ExpiresAt, 0).UTC()
	if time.Now().UTC().After(expiresAt) {
		return port.TokenClaims{}, ErrExpiredToken
	}

	return port.TokenClaims{
		UserID:    payload.UserID,
		Username:  payload.Username,
		Role:      payload.Role,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *JWTTokenService) sign(signingInput string) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(signingInput))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func encodeSegment(value any) (string, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(encoded), nil
}
