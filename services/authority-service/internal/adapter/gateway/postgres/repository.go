package postgres

import (
	"database/sql"

	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/gateway/postgres/sqlcgen"
)

type Repository struct {
	db *sql.DB
	q  sqlcgen.Querier
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db, q: sqlcgen.New(db)}
}
