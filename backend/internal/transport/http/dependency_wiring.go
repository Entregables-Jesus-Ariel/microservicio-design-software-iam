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
	db           *sql.DB
	registerUser *usecase.RegisterUser
	loginUser    *usecase.LoginUser
}

func buildDependencies(cfg config.Config) (*dependencies, error) {
	db, err := postgres.Open(cfg)
	if err != nil {
		return nil, err
	}

	users := postgres.NewUserRepository(db)
	roles := postgres.NewRoleRepository(db)
	refreshTokens := postgres.NewRefreshTokenRepository(db)
	audit := postgres.NewAuditRepository(db)

	hasher := security.NewBcryptHasher(0)
	tokens := security.NewJWTTokenService(cfg.TokenSecret, cfg.TokenTTL)

	return &dependencies{
		db:           db,
		registerUser: usecase.NewRegisterUser(users, roles, hasher),
		loginUser:    usecase.NewLoginUser(users, refreshTokens, audit, hasher, tokens),
	}, nil
}

func (d *dependencies) Close() {
	_ = d.db.Close()
}
