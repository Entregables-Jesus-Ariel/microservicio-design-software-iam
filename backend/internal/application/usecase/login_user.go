// Package usecase holds the application's business workflows. Each use
// case orchestrates ports; none of them import persistence or transport.
package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"iam/internal/application/port"
	"iam/internal/domain"
)

// maxFailedAttempts is how many consecutive wrong passwords trigger a lock.
const maxFailedAttempts = 5

// lockoutDuration is how long an account stays locked once triggered.
const lockoutDuration = 15 * time.Minute

// LoginUser authenticates a user by email and password, issuing a token
// pair on success and recording every attempt for audit purposes.
type LoginUser struct {
	users         port.UserRepository
	roles         port.RoleRepository
	refreshTokens port.RefreshTokenRepository
	audit         port.AuditRepository
	hasher        port.PasswordHasher
	tokens        port.TokenService
}

// NewLoginUser builds the use case with its dependencies.
func NewLoginUser(
	users port.UserRepository,
	roles port.RoleRepository,
	refreshTokens port.RefreshTokenRepository,
	audit port.AuditRepository,
	hasher port.PasswordHasher,
	tokens port.TokenService,
) *LoginUser {
	return &LoginUser{
		users:         users,
		roles:         roles,
		refreshTokens: refreshTokens,
		audit:         audit,
		hasher:        hasher,
		tokens:        tokens,
	}
}

// LoginUserInput carries the credentials submitted on the login form.
type LoginUserInput struct {
	Email    string
	Password string
}

// LoginUserOutput is what the transport layer returns to the client.
type LoginUserOutput struct {
	AccessToken  string
	RefreshToken string
	User         domain.User
}

// Execute validates credentials, and on success mints a token pair.
// Every branch records exactly one audit row before returning.
func (uc *LoginUser) Execute(ctx context.Context, input LoginUserInput) (LoginUserOutput, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	user, err := uc.users.FindByEmail(ctx, email)
	if errors.Is(err, domain.ErrInvalidCredentials) {
		uc.recordAttempt(ctx, nil, email, port.LoginOutcomeUserNotFound)
		return LoginUserOutput{}, domain.ErrInvalidCredentials
	}
	if err != nil {
		return LoginUserOutput{}, err
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		uc.recordAttempt(ctx, &user.ID, email, port.LoginOutcomeAccountLocked)
		return LoginUserOutput{}, domain.ErrAccountLocked
	}

	if !uc.hasher.Verify(user.PasswordHash, input.Password) {
		uc.handleFailedAttempt(ctx, user)
		uc.recordAttempt(ctx, &user.ID, email, port.LoginOutcomeInvalidPassword)
		return LoginUserOutput{}, domain.ErrInvalidCredentials
	}

	if err := uc.users.ResetFailedAttempts(ctx, user.ID); err != nil {
		return LoginUserOutput{}, err
	}

	userRoles, err := uc.roles.GetUserRoles(ctx, user.ID)
	if err != nil {
		return LoginUserOutput{}, err
	}

	accessToken, err := uc.tokens.GenerateAccessToken(port.TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Roles:  userRoles,
	})
	if err != nil {
		return LoginUserOutput{}, err
	}

	refreshPlain, refreshHash, err := uc.tokens.GenerateRefreshToken()
	if err != nil {
		return LoginUserOutput{}, err
	}
	if err := uc.refreshTokens.Store(ctx, user.ID, refreshHash, time.Now().Add(7*24*time.Hour)); err != nil {
		return LoginUserOutput{}, err
	}

	uc.recordAttempt(ctx, &user.ID, email, port.LoginOutcomeSuccess)

	return LoginUserOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshPlain,
		User:         user,
	}, nil
}

// handleFailedAttempt increments the counter and, once the threshold is
// reached, locks the account for lockoutDuration. Any error here is
// swallowed on purpose: a bookkeeping failure must not change the
// authentication outcome the caller already decided (invalid credentials).
func (uc *LoginUser) handleFailedAttempt(ctx context.Context, user domain.User) {
	var lockedUntil *time.Time
	if int(user.FailedAttempts)+1 >= maxFailedAttempts {
		until := time.Now().Add(lockoutDuration)
		lockedUntil = &until
	}
	_ = uc.users.RegisterFailedAttempt(ctx, user.ID, lockedUntil)
}

// recordAttempt swallows audit-write errors: a failed audit insert must
// never mask the real authentication result to the caller.
func (uc *LoginUser) recordAttempt(ctx context.Context, userID *string, email string, outcome port.LoginOutcome) {
	_ = uc.audit.RecordLoginAttempt(ctx, port.LoginAttempt{
		UserID:         userID,
		EmailAttempted: email,
		Outcome:        outcome,
	})
}
