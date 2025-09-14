package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/postgres/sqlcgen"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"github.com/google/uuid"
)

// keywordRepository implements gateway.KeywordRepository interface
type keywordRepository struct {
	*Repository
}

// NewKeywordRepository creates a new keyword repository
func NewKeywordRepository(repo *Repository) gateway.KeywordRepository {
	return &keywordRepository{Repository: repo}
}

// Save creates a new keyword
func (r *keywordRepository) Save(ctx context.Context, k *domain.Keyword) error {
	id, err := uuid.Parse(string(k.ID))
	if err != nil {
		return err
	}

	genreID, err := uuid.Parse(string(k.GenreID))
	if err != nil {
		return err
	}

	return r.q.CreateKeyword(ctx, sqlcgen.CreateKeywordParams{
		ID:          id,
		GenreID:     genreID,
		Name:        k.Name,
		FilterType:  string(k.FilterType),
		Pattern:     k.Pattern,
		TargetField: k.TargetField,
		Enabled:     sql.NullBool{Bool: k.Enabled, Valid: true},
		Description: toNullString(k.Description),
		CreatedAt:   sql.NullTime{Time: k.CreatedAt, Valid: true},
	})
}

// FindAll finds all keywords
func (r *keywordRepository) FindAll(ctx context.Context, enabledOnly bool) ([]*domain.Keyword, error) {
	var keywords []*domain.Keyword

	if enabledOnly {
		rows, err := r.q.ListEnabledKeywords(ctx)
		if err != nil {
			return nil, err
		}

		keywords = make([]*domain.Keyword, len(rows))
		for i, row := range rows {
			keywords[i] = &domain.Keyword{
				ID:          valueobject.UUID(row.ID.String()),
				GenreID:     valueobject.UUID(row.GenreID.String()),
				Name:        row.Name,
				FilterType:  valueobject.FilterType(row.FilterType),
				Pattern:     row.Pattern,
				TargetField: row.TargetField,
				Enabled:     row.Enabled.Bool,
				Description: nullStringToPtr(row.Description),
				CreatedAt:   row.CreatedAt.Time,
				UpdatedAt:   nullTimeToPtr(row.UpdatedAt),
			}
		}
	} else {
		// TODO: Implement ListAllKeywords query if needed
		return nil, nil
	}

	return keywords, nil
}

// FindByID finds a keyword by ID
func (r *keywordRepository) FindByID(ctx context.Context, id valueobject.UUID) (*domain.Keyword, error) {
	uid, err := uuid.Parse(string(id))
	if err != nil {
		return nil, err
	}

	row, err := r.q.GetKeywordByID(ctx, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &domain.Keyword{
		ID:          valueobject.UUID(row.ID.String()),
		GenreID:     valueobject.UUID(row.GenreID.String()),
		Name:        row.Name,
		FilterType:  valueobject.FilterType(row.FilterType),
		Pattern:     row.Pattern,
		TargetField: row.TargetField,
		Enabled:     row.Enabled.Bool,
		Description: nullStringToPtr(row.Description),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   nullTimeToPtr(row.UpdatedAt),
	}, nil
}

// Update updates an existing keyword
func (r *keywordRepository) Update(ctx context.Context, k *domain.Keyword) error {
	id, err := uuid.Parse(string(k.ID))
	if err != nil {
		return err
	}

	return r.q.UpdateKeyword(ctx, sqlcgen.UpdateKeywordParams{
		ID:          id,
		Name:        k.Name,
		FilterType:  string(k.FilterType),
		Pattern:     k.Pattern,
		TargetField: k.TargetField,
		Enabled:     sql.NullBool{Bool: k.Enabled, Valid: true},
		Description: toNullString(k.Description),
		UpdatedAt:   sql.NullTime{Time: *k.UpdatedAt, Valid: true},
	})
}

// FindByGenre finds keywords by genre ID
func (r *keywordRepository) FindByGenre(ctx context.Context, genreID valueobject.UUID, enabledOnly bool) ([]*domain.Keyword, error) {
	gid, err := uuid.Parse(string(genreID))
	if err != nil {
		return nil, err
	}

	rows, err := r.q.ListKeywordsByGenre(ctx, sqlcgen.ListKeywordsByGenreParams{
		GenreID: gid,
		Column2: enabledOnly,
	})
	if err != nil {
		return nil, err
	}

	keywords := make([]*domain.Keyword, len(rows))
	for i, row := range rows {
		keywords[i] = &domain.Keyword{
			ID:          valueobject.UUID(row.ID.String()),
			GenreID:     valueobject.UUID(row.GenreID.String()),
			Name:        row.Name,
			FilterType:  valueobject.FilterType(row.FilterType),
			Pattern:     row.Pattern,
			TargetField: row.TargetField,
			Enabled:     row.Enabled.Bool,
			Description: nullStringToPtr(row.Description),
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   nullTimeToPtr(row.UpdatedAt),
		}
	}

	return keywords, nil
}

// FindByGenreAndType finds keywords by genre ID and filter type
func (r *keywordRepository) FindByGenreAndType(ctx context.Context, genreID valueobject.UUID, filterType valueobject.FilterType, enabledOnly bool) ([]*domain.Keyword, error) {
	gid, err := uuid.Parse(string(genreID))
	if err != nil {
		return nil, err
	}

	rows, err := r.q.ListKeywordsByGenreAndType(ctx, sqlcgen.ListKeywordsByGenreAndTypeParams{
		GenreID:    gid,
		FilterType: string(filterType),
		Column3:    enabledOnly,
	})
	if err != nil {
		return nil, err
	}

	keywords := make([]*domain.Keyword, len(rows))
	for i, row := range rows {
		keywords[i] = &domain.Keyword{
			ID:          valueobject.UUID(row.ID.String()),
			GenreID:     valueobject.UUID(row.GenreID.String()),
			Name:        row.Name,
			FilterType:  valueobject.FilterType(row.FilterType),
			Pattern:     row.Pattern,
			TargetField: row.TargetField,
			Enabled:     row.Enabled.Bool,
			Description: nullStringToPtr(row.Description),
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   nullTimeToPtr(row.UpdatedAt),
		}
	}

	return keywords, nil
}

// SoftDelete soft deletes a keyword
func (r *keywordRepository) SoftDelete(ctx context.Context, id valueobject.UUID) error {
	kid, err := uuid.Parse(string(id))
	if err != nil {
		return err
	}

	now := time.Now()
	return r.q.SoftDeleteKeyword(ctx, sqlcgen.SoftDeleteKeywordParams{
		ID:        kid,
		DeletedAt: sql.NullTime{Time: now, Valid: true},
	})
}

// toNullString converts *string to sql.NullString
func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// nullStringToPtr converts sql.NullString to *string
func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// toDomainKeyword converts database row to domain keyword
func toDomainKeyword(row sqlcgen.IngestionKeyword) *domain.Keyword {
	k := &domain.Keyword{
		ID:          valueobject.UUID(row.ID.String()),
		GenreID:     valueobject.UUID(row.GenreID.String()),
		Name:        row.Name,
		FilterType:  valueobject.FilterType(row.FilterType),
		Pattern:     row.Pattern,
		TargetField: row.TargetField,
		Enabled:     row.Enabled.Bool,
		Description: nullStringToPtr(row.Description),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   nullTimeToPtr(row.UpdatedAt),
	}

	if row.DeletedAt.Valid {
		k.DeletedAt = &row.DeletedAt.Time
	}

	return k
}
