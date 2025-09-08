//go:build !sqlc

package datastore

import (
    "database/sql"
    "errors"
    outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
)

// NewRepositories is unavailable without the sqlc build tag.
func NewRepositories(db *sql.DB) (outgateway.AccountRepository, outgateway.IdentityRepository, outgateway.RoleRepository, error) {
    return nil, nil, nil, errors.New("sqlc repositories not available; build with -tags sqlc")
}

