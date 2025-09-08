package gateway

import (
    "context"
    "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain"
)

// IdentityRepository manages identities linked to accounts.
type IdentityRepository interface {
    // ListIdentitiesByAccount lists identities for an account.
    ListIdentitiesByAccount(ctx context.Context, accountID string) ([]domain.Identity, error)
    // FindByProvider finds the owning account by provider and provider uid.
    FindByProvider(ctx context.Context, provider domain.Provider, providerUID string) (*domain.Account, error)
    // Upsert creates or updates an identity for an account.
    Upsert(ctx context.Context, accountID string, id domain.Identity) error
    // Delete removes an identity from an account.
    Delete(ctx context.Context, accountID string, provider domain.Provider) error
}

