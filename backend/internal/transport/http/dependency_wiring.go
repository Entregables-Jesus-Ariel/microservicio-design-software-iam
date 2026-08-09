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
}

func buildDependencies(cfg config.Config) (*dependencies, error) {
	db, err := postgres.Open(cfg)
	if err != nil {
		return nil, err
	}

	users := postgres.NewUserRepository(db)
	roles := postgres.NewRoleRepository(db)
	hasher := security.NewBcryptHasher(0)

	return &dependencies{
		db:           db,
		registerUser: usecase.NewRegisterUser(users, roles, hasher),
	}, nil
}

func (d *dependencies) Close() {
	_ = d.db.Close()
}
