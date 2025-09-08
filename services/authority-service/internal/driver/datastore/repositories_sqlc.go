//go:build sqlc

package datastore

import (
    "database/sql"
    pgrepo "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/gateway/postgres"
    outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
)

// NewRepositories returns postgres-backed repositories generated with sqlc.
func NewRepositories(db *sql.DB) (outgateway.AccountRepository, outgateway.IdentityRepository, outgateway.RoleRepository, error) {
    r := pgrepo.New(db)
    return r, r, r, nil
}

