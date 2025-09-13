package postgres

import (
	"context"
	"database/sql"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

// transactionManager implements gateway.TransactionManager
type transactionManager struct {
	db *sql.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *sql.DB) gateway.TransactionManager {
	return &transactionManager{db: db}
}

// Execute executes a function within a database transaction
func (tm *transactionManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Store transaction in context
	txCtx := context.WithValue(ctx, txKey{}, tx)

	// Execute the function
	err = fn(txCtx)
	if err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	// Commit on success
	return tx.Commit()
}

// txKey is used as a key for storing transaction in context
type txKey struct{}

// GetTx retrieves a transaction from context
func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}