package usecase

import (
	"context"
	"time"

	"iam/internal/application/port"
)

type RefreshSession struct {
	users         port.UserRepository
	roles         port.RoleRepository
	refreshTokens port.RefreshTokenRepository
	tokens        port.TokenService
}

func NewRefreshSession(
	users port.UserRepository,
	roles port.RoleRepository,
	refreshTokens port.RefreshTokenRepository,
	tokens port.TokenService,
) *RefreshSession {
	return &RefreshSession{
		users:         users,
		roles:         roles,
		refreshTokens: refreshTokens,
		tokens:        tokens,
	}
}

type RefreshSessionInput struct {
	RefreshToken string
}

type RefreshSessionOutput struct {
	AccessToken  string
	RefreshToken string
}

func (uc *RefreshSession) Execute(ctx context.Context, input RefreshSessionInput) (RefreshSessionOutput, error) {
	tokenHash := uc.tokens.HashToken(input.RefreshToken)

	userID, err := uc.refreshTokens.Find(ctx, tokenHash)
	if err != nil {
		return RefreshSessionOutput{}, err
	}

	user, err := uc.users.FindByID(ctx, userID)
	if err != nil {
		return RefreshSessionOutput{}, err
	}

	if err := uc.refreshTokens.Revoke(ctx, tokenHash); err != nil {
		return RefreshSessionOutput{}, err
	}

	userRoles, err := uc.roles.GetUserRoles(ctx, user.ID)
	if err != nil {
		return RefreshSessionOutput{}, err
	}

	accessToken, err := uc.tokens.GenerateAccessToken(port.TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Roles:  userRoles,
	})
	if err != nil {
		return RefreshSessionOutput{}, err
	}

	refreshPlain, refreshHash, err := uc.tokens.GenerateRefreshToken()
	if err != nil {
		return RefreshSessionOutput{}, err
	}

	if err := uc.refreshTokens.Store(ctx, user.ID, refreshHash, time.Now().Add(7*24*time.Hour)); err != nil {
		return RefreshSessionOutput{}, err
	}

	return RefreshSessionOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshPlain,
	}, nil
}
