package postgres

import (
    "database/sql"

    "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/gateway/postgres/sqlcgen"
    outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
)

type Repository struct {
    db *sql.DB
    q  sqlcgen.Querier
}

func New(db *sql.DB) *Repository {
    return &Repository{db: db, q: sqlcgen.New(db)}
}

// Ensure Repository implements the output ports.
var _ outgateway.AccountRepository = (*Repository)(nil)
var _ outgateway.IdentityRepository = (*Repository)(nil)
var _ outgateway.RoleRepository = (*Repository)(nil)
