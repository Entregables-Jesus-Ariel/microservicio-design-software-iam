// Package mysql implements the application repository ports on top of
// MySQL. It owns SQL and driver details so no other layer does.
package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"ferreteria/internal/config"
)

const (
	maxOpenConnections = 25
	maxIdleConnections = 5
	connectionLifetime = 30 * time.Minute
)

// Open builds a connection pool from validated configuration. Credentials
// come from config, never from literals in this package.
//
// clientFoundRows=true makes UPDATE report matched rows (not changed rows)
// in RowsAffected — sin esto, un update legítimo cuyo valor nuevo coincide
// con el viejo se leería erróneamente como "movement not found".
func Open(settings config.Config) (*sql.DB, error) {
	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true&clientFoundRows=true&loc=UTC",
		settings.DatabaseUser,
		settings.DatabasePassword,
		settings.DatabaseHost,
		settings.DatabasePort,
		settings.DatabaseName,
	)

	database, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, fmt.Errorf("open database connection: %w", err)
	}

	database.SetMaxOpenConns(maxOpenConnections)
	database.SetMaxIdleConns(maxIdleConnections)
	database.SetConnMaxLifetime(connectionLifetime)

	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("verify database connection: %w", err)
	}
	return database, nil
}