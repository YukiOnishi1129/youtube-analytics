//go:build sqlc

package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/gateway/postgres/sqlcgen"
	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain"
	outgateway "github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/port/output/gateway"
	"github.com/google/uuid"
)

// Ensure interfaces are implemented
var (
	_ outgateway.AccountRepository  = (*Repository)(nil)
	_ outgateway.IdentityRepository = (*Repository)(nil)
	_ outgateway.RoleRepository     = (*Repository)(nil)
)

type Repository struct {
	db *sql.DB
	q  sqlcgen.Querier
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db, q: sqlcgen.New(db)}
}

// AccountRepository
func (r *Repository) FindByID(ctx context.Context, id string) (*domain.Account, error) {
	accUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	a, err := r.q.GetAccountByID(ctx, accUUID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return mapAccountByIDRow(a), nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*domain.Account, error) {
	a, err := r.q.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return mapAccountByEmailRow(a), nil
}

func (r *Repository) Save(ctx context.Context, a *domain.Account) error {
	id := uuid.Nil
	if a.ID != "" {
		if parsed, err := uuid.Parse(a.ID); err == nil {
			id = parsed
		}
	}
	if id == uuid.Nil {
		id = uuid.New()
	}
	err := r.q.UpsertAccount(ctx, sqlcgen.UpsertAccountParams{
		ID:            id,
		Email:         a.Email,
		EmailVerified: a.EmailVerified,
		DisplayName:   a.DisplayName,
		PhotoUrl:      a.PhotoURL,
		IsActive:      a.IsActive,
		LastLoginAt:   toNullTime(a.LastLoginAt),
	})
	if err != nil {
		return err
	}
	a.ID = id.String()
	return nil
}

// IdentityRepository
func (r *Repository) ListByAccount(ctx context.Context, accountID string) ([]domain.Identity, error) {
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

// Delete implements IdentityRepository.Delete
func (r *Repository) Delete(ctx context.Context, accountID string, provider domain.Provider) error {
	acc, err := uuid.Parse(accountID)
	if err != nil {
		return err
	}
	return r.q.DeleteIdentity(ctx, sqlcgen.DeleteIdentityParams{AccountID: acc, Provider: string(provider)})
}

// RoleRepository
func (r *Repository) ListByAccount(ctx context.Context, accountID string) ([]domain.Role, error) {
	acc, err := uuid.Parse(accountID)
	if err != nil {
		return nil, nil
	}
	rows, err := r.q.ListRolesByAccount(ctx, acc)
	if err != nil {
		return nil, err
	}
	roles := make([]domain.Role, 0, len(rows))
	for _, name := range rows {
		roles = append(roles, domain.Role{Name: name})
	}
	return roles, nil
}

func (r *Repository) Assign(ctx context.Context, accountID string, role domain.Role) error {
	acc, err := uuid.Parse(accountID)
	if err != nil {
		return err
	}
	return r.q.AssignRole(ctx, sqlcgen.AssignRoleParams{AccountID: acc, Name: role.Name})
}

func (r *Repository) Revoke(ctx context.Context, accountID string, role domain.Role) error {
	acc, err := uuid.Parse(accountID)
	if err != nil {
		return err
	}
	return r.q.RevokeRole(ctx, sqlcgen.RevokeRoleParams{AccountID: acc, Name: role.Name})
}

// helpers
func mapAccountByIDRow(a sqlcgen.GetAccountByIDRow) *domain.Account {
	return &domain.Account{
		ID:            a.ID.String(),
		Email:         a.Email,
		EmailVerified: a.EmailVerified,
		DisplayName:   a.DisplayName.String,
		PhotoURL:      a.PhotoUrl.String,
		IsActive:      a.IsActive,
		LastLoginAt:   fromNullTime(a.LastLoginAt),
	}
}

func mapAccountByEmailRow(a sqlcgen.GetAccountByEmailRow) *domain.Account {
	return &domain.Account{
		ID:            a.ID.String(),
		Email:         a.Email,
		EmailVerified: a.EmailVerified,
		DisplayName:   a.DisplayName.String,
		PhotoURL:      a.PhotoUrl.String,
		IsActive:      a.IsActive,
		LastLoginAt:   fromNullTime(a.LastLoginAt),
	}
}

func mapAccountByProviderRow(a sqlcgen.FindAccountByProviderRow) *domain.Account {
	return &domain.Account{
		ID:            a.ID.String(),
		Email:         a.Email,
		EmailVerified: a.EmailVerified,
		DisplayName:   a.DisplayName.String,
		PhotoURL:      a.PhotoUrl.String,
		IsActive:      a.IsActive,
		LastLoginAt:   fromNullTime(a.LastLoginAt),
	}
}

func toNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func fromNullTime(t sql.NullTime) *time.Time {
	if !t.Valid {
		return nil
	}
	tm := t.Time
	return &tm
}

func sqlNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
