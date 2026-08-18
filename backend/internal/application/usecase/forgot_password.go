package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"iam/internal/application/port"
	"iam/internal/domain"
)

type ForgotPassword struct {
	users        port.UserRepository
	resetTokens  port.PasswordResetRepository
	tokenService port.TokenService
}

func NewForgotPassword(
	users port.UserRepository,
	resetTokens port.PasswordResetRepository,
	tokenService port.TokenService,
) *ForgotPassword {
	return &ForgotPassword{
		users:        users,
		resetTokens:  resetTokens,
		tokenService: tokenService,
	}
}

type ForgotPasswordInput struct {
	Email string
}

type ForgotPasswordOutput struct {
	// We return the plain token so it can be returned by the HTTP handler
	// since we don't have an email service yet.
	ResetToken string
}

func (uc *ForgotPassword) Execute(ctx context.Context, input ForgotPasswordInput) (ForgotPasswordOutput, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	user, err := uc.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			// Do not leak whether the email exists or not.
			// Return a generic success-like state or silently drop.
			// But for dev, we can just return a fake success.
			return ForgotPasswordOutput{}, nil
		}
		return ForgotPasswordOutput{}, err
	}

	plainToken, hash, err := uc.tokenService.GenerateRefreshToken() // We can reuse this method to generate a secure random token
	if err != nil {
		return ForgotPasswordOutput{}, err
	}

	expiresAt := time.Now().Add(1 * time.Hour)
	if err := uc.resetTokens.Store(ctx, user.ID, hash, expiresAt); err != nil {
		return ForgotPasswordOutput{}, err
	}

	return ForgotPasswordOutput{
		ResetToken: plainToken,
	}, nil
}
