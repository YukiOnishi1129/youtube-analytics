package postgres

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/gateway/postgres/sqlcgen"
	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain"
	"github.com/google/uuid"
)

// IdentityRepository methods
func (r *Repository) ListIdentitiesByAccount(ctx context.Context, accountID string) ([]domain.Identity, error) {
	id, err := uuid.Parse(accountID)
	if err != nil {
		return nil, nil
	}
	rows, err := r.q.ListIdentitiesByAccount(ctx, id)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Identity, 0, len(rows))
	for _, row := range rows {
		res = append(res, domain.Identity{Provider: domain.Provider(row.Provider), ProviderUID: row.ProviderUid.String})
	}
	return res, nil
}

func (r *Repository) FindByProvider(ctx context.Context, provider domain.Provider, providerUID string) (*domain.Account, error) {
	a, err := r.q.FindAccountByProvider(ctx, sqlcgen.FindAccountByProviderParams{Provider: string(provider), ProviderUid: sqlNullString(providerUID)})
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return mapAccountByProviderRow(a), nil
}

func (r *Repository) Upsert(ctx context.Context, accountID string, id domain.Identity) error {
	acc, err := uuid.Parse(accountID)
	if err != nil {
		return err
	}
	return r.q.UpsertIdentity(ctx, sqlcgen.UpsertIdentityParams{ID: uuid.New(), AccountID: acc, Provider: string(id.Provider), ProviderUid: sqlNullString(id.ProviderUID)})
}

func (r *Repository) Delete(ctx context.Context, accountID string, provider domain.Provider) error {
	acc, err := uuid.Parse(accountID)
	if err != nil {
		return err
	}
	return r.q.DeleteIdentity(ctx, sqlcgen.DeleteIdentityParams{AccountID: acc, Provider: string(provider)})
}
