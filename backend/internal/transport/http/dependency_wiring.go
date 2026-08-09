package http

import (
	"database/sql"

	"ferreteria/internal/application/usecase"
	"ferreteria/internal/config"
	"ferreteria/internal/persistence/mysql"
	"ferreteria/internal/security"
)

// dependencies holds everything the handlers need, built once at boot.
type dependencies struct {
	database         *sql.DB
	tokens           *security.JWTTokenService
	recordMovement   *usecase.RecordMovement
	listMovement     *usecase.ListMovement
	editMovement     *usecase.EditMovement
	cancelMovement   *usecase.CancelMovement
	getMovementAudit *usecase.GetMovementAudit
	getPeriodSummary *usecase.GetPeriodSummary
	listCategory     *usecase.ListCategory
	createCategory   *usecase.CreateCategory
	authenticateUser *usecase.AuthenticateUser
}

// buildDependencies opens persistence and composes the use cases.
func buildDependencies(settings config.Config) (*dependencies, error) {
	database, err := mysql.Open(settings)
	if err != nil {
		return nil, err
	}

	movements := mysql.NewMovementRepository(database)
	categories := mysql.NewCategoryRepository(database)
	users := mysql.NewUserRepository(database)
	audits := mysql.NewMovementAuditRepository(database)

	hasher := security.NewBcryptHasher(0)
	tokens := security.NewJWTTokenService(settings.TokenSecret)

	return &dependencies{
		database:         database,
		tokens:           tokens,
		recordMovement:   usecase.NewRecordMovement(movements, categories, audits),
		listMovement:     usecase.NewListMovement(movements),
		editMovement:     usecase.NewEditMovement(movements, categories, audits),
		cancelMovement:   usecase.NewCancelMovement(movements, audits),
		getMovementAudit: usecase.NewGetMovementAudit(movements, audits),
		getPeriodSummary: usecase.NewGetPeriodSummary(movements),
		listCategory:     usecase.NewListCategory(categories),
		createCategory:   usecase.NewCreateCategory(categories),
		authenticateUser: usecase.NewAuthenticateUser(users, hasher, tokens, settings.TokenTTL),
	}, nil
}

// Close releases the persistence handle.
func (d *dependencies) Close() {
	if d.database != nil {
		_ = d.database.Close()
	}
}