package repository

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// KeywordGroupRepository interface for keyword group persistence
type KeywordGroupRepository interface {
	// Create creates a new keyword group with its items
	Create(ctx context.Context, group *domain.KeywordGroup) error

	// Update updates an existing keyword group (excluding items)
	Update(ctx context.Context, group *domain.KeywordGroup) error

	// UpdateWithItems updates a keyword group and replaces all its items
	UpdateWithItems(ctx context.Context, group *domain.KeywordGroup) error

	// Delete soft deletes a keyword group
	Delete(ctx context.Context, id valueobject.UUID) error

	// FindByID finds a keyword group by ID including its items
	FindByID(ctx context.Context, id valueobject.UUID) (*domain.KeywordGroup, error)

	// FindByGenreID finds all keyword groups for a genre
	FindByGenreID(ctx context.Context, genreID valueobject.UUID) ([]*domain.KeywordGroup, error)

	// List lists keyword groups with pagination
	List(ctx context.Context, limit, offset int) ([]*domain.KeywordGroup, error)

	// ListByEnabled lists enabled keyword groups
	ListByEnabled(ctx context.Context, enabled bool, limit, offset int) ([]*domain.KeywordGroup, error)
}