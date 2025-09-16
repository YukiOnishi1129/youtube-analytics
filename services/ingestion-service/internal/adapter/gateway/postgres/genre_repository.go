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

// genreRepository implements gateway.GenreRepository interface
type genreRepository struct {
	*Repository
}

// NewGenreRepository creates a new genre repository
func NewGenreRepository(repo *Repository) gateway.GenreRepository {
	return &genreRepository{Repository: repo}
}

// Save creates a new genre
func (r *genreRepository) Save(ctx context.Context, g *domain.Genre) error {
	id, err := uuid.Parse(string(g.ID))
	if err != nil {
		return err
	}

	categoryIDs := make([]int32, len(g.CategoryIDs))
	for i, catID := range g.CategoryIDs {
		categoryIDs[i] = int32(catID)
	}

	now := time.Now()
	return r.q.CreateGenre(ctx, sqlcgen.CreateGenreParams{
		ID:          id,
		Code:        g.Code,
		Name:        g.Name,
		Language:    g.Language,
		RegionCode:  g.RegionCode,
		CategoryIds: categoryIDs,
		Enabled:     g.Enabled,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}

// Update updates an existing genre
func (r *genreRepository) Update(ctx context.Context, g *domain.Genre) error {
	id, err := uuid.Parse(string(g.ID))
	if err != nil {
		return err
	}

	categoryIDs := make([]int32, len(g.CategoryIDs))
	for i, catID := range g.CategoryIDs {
		categoryIDs[i] = int32(catID)
	}

	return r.q.UpdateGenre(ctx, sqlcgen.UpdateGenreParams{
		ID:          id,
		Name:        g.Name,
		CategoryIds: categoryIDs,
		Enabled:     g.Enabled,
		UpdatedAt:   time.Now(),
	})
}

// FindByID finds a genre by ID
func (r *genreRepository) FindByID(ctx context.Context, id valueobject.UUID) (*domain.Genre, error) {
	genreID, err := uuid.Parse(string(id))
	if err != nil {
		return nil, err
	}

	row, err := r.q.GetGenreByID(ctx, genreID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainGenre(row), nil
}

// FindByCode finds a genre by code
func (r *genreRepository) FindByCode(ctx context.Context, code string) (*domain.Genre, error) {
	row, err := r.q.GetGenreByCode(ctx, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainGenre(row), nil
}

// FindAll finds all genres
func (r *genreRepository) FindAll(ctx context.Context) ([]*domain.Genre, error) {
	rows, err := r.q.ListGenres(ctx)
	if err != nil {
		return nil, err
	}

	genres := make([]*domain.Genre, len(rows))
	for i, row := range rows {
		genres[i] = toDomainGenre(row)
	}

	return genres, nil
}

// FindEnabled finds all enabled genres
func (r *genreRepository) FindEnabled(ctx context.Context) ([]*domain.Genre, error) {
	rows, err := r.q.ListEnabledGenres(ctx)
	if err != nil {
		return nil, err
	}

	genres := make([]*domain.Genre, len(rows))
	for i, row := range rows {
		genres[i] = toDomainGenre(row)
	}

	return genres, nil
}

// toDomainGenre converts a database row to a domain genre
func toDomainGenre(row sqlcgen.IngestionGenre) *domain.Genre {
	categoryIDs := make([]valueobject.CategoryID, len(row.CategoryIds))
	for i, id := range row.CategoryIds {
		categoryIDs[i] = valueobject.CategoryID(id)
	}

	return &domain.Genre{
		ID:          valueobject.UUID(row.ID.String()),
		Code:        row.Code,
		Name:        row.Name,
		Language:    row.Language,
		RegionCode:  row.RegionCode,
		CategoryIDs: categoryIDs,
		Enabled:     row.Enabled,
	}
}