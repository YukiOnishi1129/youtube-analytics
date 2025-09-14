package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/postgres/sqlcgen"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

// youtubeCategoryRepository implements gateway.YouTubeCategoryRepository interface
type youtubeCategoryRepository struct {
	*Repository
}

// NewYouTubeCategoryRepository creates a new YouTube category repository
func NewYouTubeCategoryRepository(repo *Repository) gateway.YouTubeCategoryRepository {
	return &youtubeCategoryRepository{Repository: repo}
}

// Save creates a new YouTube category
func (r *youtubeCategoryRepository) Save(ctx context.Context, c *domain.YouTubeCategory) error {
	now := time.Now()
	return r.q.CreateYouTubeCategory(ctx, sqlcgen.CreateYouTubeCategoryParams{
		ID:         int32(c.ID),
		Name:       c.Name,
		Assignable: c.Assignable,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
}

// Update updates an existing YouTube category
func (r *youtubeCategoryRepository) Update(ctx context.Context, c *domain.YouTubeCategory) error {
	return r.q.UpdateYouTubeCategory(ctx, sqlcgen.UpdateYouTubeCategoryParams{
		ID:         int32(c.ID),
		Name:       c.Name,
		Assignable: c.Assignable,
		UpdatedAt:  time.Now(),
	})
}

// FindByID finds a YouTube category by ID
func (r *youtubeCategoryRepository) FindByID(ctx context.Context, id valueobject.CategoryID) (*domain.YouTubeCategory, error) {
	row, err := r.q.GetYouTubeCategoryByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainYouTubeCategory(row), nil
}

// FindAll finds all YouTube categories
func (r *youtubeCategoryRepository) FindAll(ctx context.Context) ([]*domain.YouTubeCategory, error) {
	rows, err := r.q.ListYouTubeCategories(ctx)
	if err != nil {
		return nil, err
	}

	categories := make([]*domain.YouTubeCategory, len(rows))
	for i, row := range rows {
		categories[i] = toDomainYouTubeCategory(row)
	}

	return categories, nil
}

// FindAssignable finds all assignable YouTube categories
func (r *youtubeCategoryRepository) FindAssignable(ctx context.Context) ([]*domain.YouTubeCategory, error) {
	rows, err := r.q.ListAssignableYouTubeCategories(ctx)
	if err != nil {
		return nil, err
	}

	categories := make([]*domain.YouTubeCategory, len(rows))
	for i, row := range rows {
		categories[i] = toDomainYouTubeCategory(row)
	}

	return categories, nil
}

// toDomainYouTubeCategory converts a database row to a domain YouTube category
func toDomainYouTubeCategory(row sqlcgen.IngestionYoutubeCategory) *domain.YouTubeCategory {
	return &domain.YouTubeCategory{
		ID:         valueobject.CategoryID(row.ID),
		Name:       row.Name,
		Assignable: row.Assignable,
	}
}