package gateway

import (
    "context"
    account "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain/account"
)

// AccountRepository is an output port owned by the use case layer.
// Implementations live in the adapter/gateway layer (e.g., Postgres, in-memory).
type AccountRepository interface {
    // GetByID returns the account with identities or ErrNotFound.
    GetByID(ctx context.Context, id string) (*account.Account, error)

    // GetByEmail returns the account by email (case-insensitive) or ErrNotFound.
    GetByEmail(ctx context.Context, email string) (*account.Account, error)

    // Create inserts a new account with identities.
    Create(ctx context.Context, a *account.Account) error

    // Update updates an existing account and its identities as needed.
    Update(ctx context.Context, a *account.Account) error

    // LinkProvider links a provider to an existing account.
    LinkProvider(ctx context.Context, accountID string, provider account.Provider, providerUID string) error
}

