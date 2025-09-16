package postgres

import (
	"context"
	"database/sql"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/postgres/sqlcgen"
)

// Repository holds the database connection and sqlc querier
type Repository struct {
	db *sql.DB
	q  sqlcgen.Querier
}

// NewRepository creates a new repository instance
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
		q:  sqlcgen.New(db),
	}
}

// WithTx returns a new repository with a transaction
func (r *Repository) WithTx(tx *sql.Tx) *Repository {
	return &Repository{
		db: r.db,
		q:  sqlcgen.New(tx),
	}
}

// DB returns the database connection
func (r *Repository) DB() *sql.DB {
	return r.db
}

// ExecTx executes a function within a database transaction
func (r *Repository) ExecTx(ctx context.Context, fn func(*Repository) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	repo := r.WithTx(tx)
	err = fn(repo)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()
}