package http

import (
	"database/sql"

	"iam/internal/application/usecase"
	"iam/internal/config"
	"iam/internal/persistence/postgres"
	"iam/internal/security"
)

// dependencies wires ports to their concrete adapters, once, at startup.
type dependencies struct {
	db             *sql.DB
	registerUser   *usecase.RegisterUser
	loginUser      *usecase.LoginUser
	refreshSession *usecase.RefreshSession
	logoutUser     *usecase.LogoutUser
	forgotPassword *usecase.ForgotPassword
	resetPassword  *usecase.ResetPassword
	rbacUsecases   *usecase.RBACUsecases
	auditUsecases  *usecase.AuditUsecases
}

func buildDependencies(cfg config.Config) (*dependencies, error) {
	db, err := postgres.Open(cfg)
	if err != nil {
		return nil, err
	}

	users := postgres.NewUserRepository(db)
	roles := postgres.NewRoleRepository(db)
	refreshTokens := postgres.NewRefreshTokenRepository(db)
	resetTokens := postgres.NewPasswordResetRepository(db)
	audit := postgres.NewAuditRepository(db)

	hasher := security.NewBcryptHasher(0)
	tokens := security.NewJWTTokenService(cfg.TokenSecret, cfg.TokenTTL)

	return &dependencies{
		db:             db,
		registerUser:   usecase.NewRegisterUser(users, roles, hasher),
		loginUser:      usecase.NewLoginUser(users, roles, refreshTokens, audit, hasher, tokens),
		refreshSession: usecase.NewRefreshSession(users, roles, refreshTokens, tokens),
		logoutUser:     usecase.NewLogoutUser(refreshTokens, tokens),
		forgotPassword: usecase.NewForgotPassword(users, resetTokens, tokens),
		resetPassword:  usecase.NewResetPassword(users, resetTokens, hasher, tokens),
		rbacUsecases:   usecase.NewRBACUsecases(users, roles),
		auditUsecases:  usecase.NewAuditUsecases(audit),
	}, nil
}

func (d *dependencies) Close() {
	_ = d.db.Close()
}
