package usecase

import (
	"context"
	"errors"
	"time"

	"ferreteria/internal/application/port"
	"ferreteria/internal/domain"
)

// AuthenticateUser verifies credentials and issues an access token.
type AuthenticateUser struct {
	users  port.UserRepository
	hasher port.PasswordHasher
	tokens port.TokenService
	ttl    time.Duration
}

// NewAuthenticateUser wires the use case with its ports.
func NewAuthenticateUser(
	users port.UserRepository,
	hasher port.PasswordHasher,
	tokens port.TokenService,
	ttl time.Duration,
) *AuthenticateUser {
	return &AuthenticateUser{users: users, hasher: hasher, tokens: tokens, ttl: ttl}
}

// Execute returns a signed token when the credentials are valid. A missing
// user and a wrong password return the same error so the response cannot
// be used to enumerate accounts.
func (u *AuthenticateUser) Execute(
	ctx context.Context,
	username string,
	password string,
) (string, error) {
	user, err := u.users.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	if !u.hasher.Verify(user.PasswordHash, password) {
		return "", domain.ErrInvalidCredentials
	}

	return u.tokens.Issue(port.TokenClaims{
		UserID:    user.ID,
		Username:  user.Username,
		Role:      user.Role,
		ExpiresAt: time.Now().UTC().Add(u.ttl),
	})
}
