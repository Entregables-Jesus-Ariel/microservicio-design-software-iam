// Package postgres implements the application ports against Postgres
// using database/sql with the pgx stdlib driver.
package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"iam/internal/config"
)

// Open connects to Postgres and verifies the connection is alive.
func Open(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.ConnString())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}
