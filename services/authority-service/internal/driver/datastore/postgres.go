package datastore

import (
    "database/sql"
)

// OpenPostgres opens a Postgres connection using the pgx stdlib driver.
// Requires building with the 'postgres' build tag to register the driver.
func OpenPostgres(dsn string) (*sql.DB, error) {
    return sql.Open("pgx", dsn)
}

