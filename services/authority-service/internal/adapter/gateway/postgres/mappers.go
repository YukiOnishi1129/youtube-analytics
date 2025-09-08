package postgres

import (
	"database/sql"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/adapter/gateway/postgres/sqlcgen"
	"github.com/YukiOnishi1129/youtube-analytics/services/authority-service/internal/domain"
)

func mapAccountByIDRow(a sqlcgen.GetAccountByIDRow) *domain.Account {
	return &domain.Account{
		ID:            a.ID.String(),
		Email:         a.Email,
		EmailVerified: a.EmailVerified.Bool,
		DisplayName:   a.DisplayName.String,
		PhotoURL:      a.PhotoUrl.String,
		IsActive:      a.IsActive.Bool,
		LastLoginAt:   fromNullTime(a.LastLoginAt),
	}
}

func mapAccountByEmailRow(a sqlcgen.GetAccountByEmailRow) *domain.Account {
	return &domain.Account{
		ID:            a.ID.String(),
		Email:         a.Email,
		EmailVerified: a.EmailVerified.Bool,
		DisplayName:   a.DisplayName.String,
		PhotoURL:      a.PhotoUrl.String,
		IsActive:      a.IsActive.Bool,
		LastLoginAt:   fromNullTime(a.LastLoginAt),
	}
}

func mapAccountByProviderRow(a sqlcgen.FindAccountByProviderRow) *domain.Account {
	return &domain.Account{
		ID:            a.ID.String(),
		Email:         a.Email,
		EmailVerified: a.EmailVerified.Bool,
		DisplayName:   a.DisplayName.String,
		PhotoURL:      a.PhotoUrl.String,
		IsActive:      a.IsActive.Bool,
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

func sqlNullBool(b bool) sql.NullBool { return sql.NullBool{Bool: b, Valid: true} }
