package gateway

import (
    "context"
    account "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain/account"
)

// IdentityRepository manages identities linked to accounts.
type IdentityRepository interface {
    // ListByAccount lists identities for an account.
    ListByAccount(ctx context.Context, accountID string) ([]account.Identity, error)

    // FindByProvider finds the owning account by provider and provider uid.
    FindByProvider(ctx context.Context, provider account.Provider, providerUID string) (*account.Account, error)

    // Save creates or updates an identity for an account.
    Save(ctx context.Context, accountID string, id account.Identity) error

    // Delete removes an identity from an account.
    Delete(ctx context.Context, accountID string, provider account.Provider) error
}

