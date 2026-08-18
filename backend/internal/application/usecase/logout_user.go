package usecase

import (
	"context"

	"iam/internal/application/port"
)

type LogoutUser struct {
	refreshTokens port.RefreshTokenRepository
	tokens        port.TokenService
}

func NewLogoutUser(
	refreshTokens port.RefreshTokenRepository,
	tokens port.TokenService,
) *LogoutUser {
	return &LogoutUser{
		refreshTokens: refreshTokens,
		tokens:        tokens,
	}
}

type LogoutUserInput struct {
	RefreshToken string
}

func (uc *LogoutUser) Execute(ctx context.Context, input LogoutUserInput) error {
	tokenHash := uc.tokens.HashToken(input.RefreshToken)
	// We just revoke it. If it doesn't exist or is already revoked, it's fine.
	return uc.refreshTokens.Revoke(ctx, tokenHash)
}
