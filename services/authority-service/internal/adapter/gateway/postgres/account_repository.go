package postgres

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/gateway/postgres/sqlcgen"
	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain"
	"github.com/google/uuid"
)

// AccountRepository methods
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
		EmailVerified: sqlNullBool(a.EmailVerified),
		DisplayName:   sqlNullString(a.DisplayName),
		PhotoUrl:      sqlNullString(a.PhotoURL),
		IsActive:      sqlNullBool(a.IsActive),
		LastLoginAt:   toNullTime(a.LastLoginAt),
	})
	if err != nil {
		return err
	}
	a.ID = id.String()
	return nil
}
