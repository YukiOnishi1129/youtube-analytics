package gateway

import "context"

// TransactionManager manages database transactions
type TransactionManager interface {
	Execute(ctx context.Context, fn func(ctx context.Context) error) error
}