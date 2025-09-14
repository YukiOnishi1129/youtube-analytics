package input

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
)

// YouTubeCategoryInputPort is the interface for YouTube category use cases
type YouTubeCategoryInputPort interface {
	ListYouTubeCategories(ctx context.Context) ([]*domain.YouTubeCategory, error)
	ListAssignableYouTubeCategories(ctx context.Context) ([]*domain.YouTubeCategory, error)
	GetYouTubeCategory(ctx context.Context, categoryID int) (*domain.YouTubeCategory, error)
	CreateYouTubeCategory(ctx context.Context, input *CreateYouTubeCategoryInput) (*domain.YouTubeCategory, error)
	UpdateYouTubeCategory(ctx context.Context, input *UpdateYouTubeCategoryInput) (*domain.YouTubeCategory, error)
}

// CreateYouTubeCategoryInput represents the input for creating a YouTube category
type CreateYouTubeCategoryInput struct {
	ID         int
	Name       string
	Assignable bool
}

// UpdateYouTubeCategoryInput represents the input for updating a YouTube category
type UpdateYouTubeCategoryInput struct {
	CategoryID int
	Name       string
	Assignable bool
}