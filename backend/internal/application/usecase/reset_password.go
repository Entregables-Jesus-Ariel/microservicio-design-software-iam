package usecase

import (
	"context"
	"errors"
	"time"

	"iam/internal/application/port"
	"iam/internal/domain"
)

type ResetPassword struct {
	users        port.UserRepository
	resetTokens  port.PasswordResetRepository
	hasher       port.PasswordHasher
	tokenService port.TokenService
}

func NewResetPassword(
	users port.UserRepository,
	resetTokens port.PasswordResetRepository,
	hasher port.PasswordHasher,
	tokenService port.TokenService,
) *ResetPassword {
	return &ResetPassword{
		users:        users,
		resetTokens:  resetTokens,
		hasher:       hasher,
		tokenService: tokenService,
	}
}

type ResetPasswordInput struct {
	Token       string
	NewPassword string
}

func (uc *ResetPassword) Execute(ctx context.Context, input ResetPasswordInput) error {
	tokenHash := uc.tokenService.HashToken(input.Token)

	req, err := uc.resetTokens.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return err
	}

	if req.IsUsed {
		return domain.ErrInvalidToken
	}

	if req.ExpiresAt.Before(time.Now()) {
		return domain.ErrInvalidToken
	}

	newHash, err := uc.hasher.Hash(input.NewPassword)
	if err != nil {
		return err
	}

	if err := uc.users.UpdatePassword(ctx, req.UserID, newHash); err != nil {
		return err
	}

	if err := uc.resetTokens.MarkAsUsed(ctx, req.ID); err != nil {
		// Log the error but it shouldn't fail the password reset process if it's already done
		// We'll return it though.
		return err
	}

	return nil
}
