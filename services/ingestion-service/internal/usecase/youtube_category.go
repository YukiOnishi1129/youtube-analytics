package usecase

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

// youtubeCategoryUseCase implements the YouTubeCategoryInputPort interface
type youtubeCategoryUseCase struct {
	categoryRepo gateway.YouTubeCategoryRepository
}

// NewYouTubeCategoryUseCase creates a new YouTube category use case
func NewYouTubeCategoryUseCase(categoryRepo gateway.YouTubeCategoryRepository) input.YouTubeCategoryInputPort {
	return &youtubeCategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

// ListYouTubeCategories lists all YouTube categories
func (u *youtubeCategoryUseCase) ListYouTubeCategories(ctx context.Context) ([]*domain.YouTubeCategory, error) {
	return u.categoryRepo.FindAll(ctx)
}

// ListAssignableYouTubeCategories lists only assignable YouTube categories
func (u *youtubeCategoryUseCase) ListAssignableYouTubeCategories(ctx context.Context) ([]*domain.YouTubeCategory, error) {
	return u.categoryRepo.FindAssignable(ctx)
}

// GetYouTubeCategory gets a YouTube category by ID
func (u *youtubeCategoryUseCase) GetYouTubeCategory(ctx context.Context, categoryID int) (*domain.YouTubeCategory, error) {
	return u.categoryRepo.FindByID(ctx, valueobject.CategoryID(categoryID))
}

// CreateYouTubeCategory creates a new YouTube category
func (u *youtubeCategoryUseCase) CreateYouTubeCategory(ctx context.Context, input *input.CreateYouTubeCategoryInput) (*domain.YouTubeCategory, error) {
	// Create domain object
	category, err := domain.NewYouTubeCategory(
		valueobject.CategoryID(input.ID),
		input.Name,
		input.Assignable,
	)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := u.categoryRepo.Save(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateYouTubeCategory updates an existing YouTube category
func (u *youtubeCategoryUseCase) UpdateYouTubeCategory(ctx context.Context, input *input.UpdateYouTubeCategoryInput) (*domain.YouTubeCategory, error) {
	// Find existing category
	category, err := u.categoryRepo.FindByID(ctx, valueobject.CategoryID(input.CategoryID))
	if err != nil {
		return nil, err
	}

	// Update domain object
	if err := category.UpdateCategory(input.Name, input.Assignable); err != nil {
		return nil, err
	}

	// Save to repository
	if err := u.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}
