package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"

	_ "github.com/lib/pq"
)

// OpenPostgres opens a Postgres connection.
// If dsn is empty, it constructs one from environment variables:
//
//	DB_HOST, DB_PORT (default 5432), DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE (default disable)
func OpenPostgres(dsn string) (*sql.DB, error) {
	if dsn == "" {
		var err error
		dsn, err = buildDSNFromEnv()
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func buildDSNFromEnv() (string, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	ssl := os.Getenv("DB_SSLMODE")
	if ssl == "" {
		ssl = "disable"
	}

	if host == "" || user == "" || name == "" {
		return "", errors.New("database config missing: require DB_HOST, DB_USER, DB_NAME (or set DATABASE_URL)")
	}

	// postgres://user:pass@host:port/dbname?sslmode=xxx
	u := &url.URL{Scheme: "postgres", Host: fmt.Sprintf("%s:%s", host, port), Path: "/" + name}
	if user != "" {
		u.User = url.UserPassword(user, pass)
	}
	q := u.Query()
	if ssl != "" {
		q.Set("sslmode", ssl)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
