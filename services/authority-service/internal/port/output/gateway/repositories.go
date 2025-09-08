package gateway

import (
    "context"
    "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain"
)

// AccountRepository is an output port owned by the use case layer.
type AccountRepository interface {
    // FindByID returns the account or ErrNotFound.
    FindByID(ctx context.Context, id string) (*domain.Account, error)
    // FindByEmail returns the account by email (case-insensitive) or ErrNotFound.
    FindByEmail(ctx context.Context, email string) (*domain.Account, error)
    // Save creates or updates an account.
    Save(ctx context.Context, a *domain.Account) error
}

// IdentityRepository manages identities linked to accounts.
type IdentityRepository interface {
    // ListByAccount lists identities for an account.
    ListByAccount(ctx context.Context, accountID string) ([]domain.Identity, error)
    // FindByProvider finds the owning account by provider and provider uid.
    FindByProvider(ctx context.Context, provider domain.Provider, providerUID string) (*domain.Account, error)
    // Save creates or updates an identity for an account.
    Save(ctx context.Context, accountID string, id domain.Identity) error
    // Delete removes an identity from an account.
    Delete(ctx context.Context, accountID string, provider domain.Provider) error
}

// RoleRepository manages role assignment for accounts.
type RoleRepository interface {
    // ListByAccount lists roles assigned to an account.
    ListByAccount(ctx context.Context, accountID string) ([]domain.Role, error)
    // Assign assigns a role to an account.
    Assign(ctx context.Context, accountID string, role domain.Role) error
    // Revoke revokes a role from an account.
    Revoke(ctx context.Context, accountID string, role domain.Role) error
}

