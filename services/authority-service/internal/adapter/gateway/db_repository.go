package gateway

import (
    "context"
    "errors"
    "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain"
    outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
)

// InMemory repositories for scaffolding. Replace with Postgres implementations.

type InMemoryAccountRepo struct{ byEmail map[string]*domain.Account }

var _ outgateway.AccountRepository = (*InMemoryAccountRepo)(nil)

func NewInMemoryAccountRepo() *InMemoryAccountRepo { return &InMemoryAccountRepo{byEmail: map[string]*domain.Account{}} }

func (r *InMemoryAccountRepo) FindByID(ctx context.Context, id string) (*domain.Account, error) {
    for _, a := range r.byEmail {
        if a.ID == id {
            return a, nil
        }
    }
    return nil, domain.ErrNotFound
}

func (r *InMemoryAccountRepo) FindByEmail(ctx context.Context, email string) (*domain.Account, error) {
    if a, ok := r.byEmail[email]; ok {
        return a, nil
    }
    return nil, domain.ErrNotFound
}

func (r *InMemoryAccountRepo) Save(ctx context.Context, a *domain.Account) error {
    if a == nil || a.Email == "" {
        return errors.New("invalid account")
    }
    r.byEmail[a.Email] = a
    return nil
}

type InMemoryIdentityRepo struct{}

var _ outgateway.IdentityRepository = (*InMemoryIdentityRepo)(nil)

func (r *InMemoryIdentityRepo) ListByAccount(ctx context.Context, accountID string) ([]domain.Identity, error) {
    return nil, nil
}

func (r *InMemoryIdentityRepo) FindByProvider(ctx context.Context, provider domain.Provider, providerUID string) (*domain.Account, error) {
    return nil, domain.ErrNotFound
}

func (r *InMemoryIdentityRepo) Save(ctx context.Context, accountID string, id domain.Identity) error { return nil }
func (r *InMemoryIdentityRepo) Delete(ctx context.Context, accountID string, provider domain.Provider) error {
    return nil
}

type InMemoryRoleRepo struct{}

var _ outgateway.RoleRepository = (*InMemoryRoleRepo)(nil)

func (r *InMemoryRoleRepo) ListByAccount(ctx context.Context, accountID string) ([]domain.Role, error) {
    return nil, nil
}

func (r *InMemoryRoleRepo) Assign(ctx context.Context, accountID string, role domain.Role) error { return nil }
func (r *InMemoryRoleRepo) Revoke(ctx context.Context, accountID string, role domain.Role) error { return nil }

