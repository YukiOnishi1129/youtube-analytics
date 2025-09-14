package usecase

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"github.com/google/uuid"
)

// genreUseCase implements the GenreInputPort interface
type genreUseCase struct {
	genreRepo gateway.GenreRepository
}

// NewGenreUseCase creates a new genre use case
func NewGenreUseCase(genreRepo gateway.GenreRepository) input.GenreInputPort {
	return &genreUseCase{
		genreRepo: genreRepo,
	}
}

// ListGenres lists all genres
func (u *genreUseCase) ListGenres(ctx context.Context) ([]*domain.Genre, error) {
	return u.genreRepo.FindAll(ctx)
}

// ListEnabledGenres lists only enabled genres
func (u *genreUseCase) ListEnabledGenres(ctx context.Context) ([]*domain.Genre, error) {
	return u.genreRepo.FindEnabled(ctx)
}

// GetGenre gets a genre by ID
func (u *genreUseCase) GetGenre(ctx context.Context, genreID uuid.UUID) (*domain.Genre, error) {
	return u.genreRepo.FindByID(ctx, valueobject.UUID(genreID.String()))
}

// GetGenreByCode gets a genre by code
func (u *genreUseCase) GetGenreByCode(ctx context.Context, code string) (*domain.Genre, error) {
	return u.genreRepo.FindByCode(ctx, code)
}

// CreateGenre creates a new genre
func (u *genreUseCase) CreateGenre(ctx context.Context, input *input.CreateGenreInput) (*domain.Genre, error) {
	// Generate new UUID
	id := uuid.New()

	// Convert category IDs
	categoryIDs := make([]valueobject.CategoryID, len(input.CategoryIDs))
	for i, catID := range input.CategoryIDs {
		categoryIDs[i] = valueobject.CategoryID(catID)
	}

	// Create domain object
	genre, err := domain.NewGenre(
		valueobject.UUID(id.String()),
		input.Code,
		input.Name,
		input.Language,
		input.RegionCode,
		categoryIDs,
	)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := u.genreRepo.Save(ctx, genre); err != nil {
		return nil, err
	}

	return genre, nil
}

// UpdateGenre updates an existing genre
func (u *genreUseCase) UpdateGenre(ctx context.Context, input *input.UpdateGenreInput) (*domain.Genre, error) {
	// Find existing genre
	genre, err := u.genreRepo.FindByID(ctx, valueobject.UUID(input.GenreID.String()))
	if err != nil {
		return nil, err
	}

	// Convert category IDs
	categoryIDs := make([]valueobject.CategoryID, len(input.CategoryIDs))
	for i, catID := range input.CategoryIDs {
		categoryIDs[i] = valueobject.CategoryID(catID)
	}

	// Update domain object
	if err := genre.Update(input.Name, categoryIDs); err != nil {
		return nil, err
	}

	// Save to repository
	if err := u.genreRepo.Update(ctx, genre); err != nil {
		return nil, err
	}

	return genre, nil
}

// EnableGenre enables a genre
func (u *genreUseCase) EnableGenre(ctx context.Context, genreID uuid.UUID) (*domain.Genre, error) {
	// Find existing genre
	genre, err := u.genreRepo.FindByID(ctx, valueobject.UUID(genreID.String()))
	if err != nil {
		return nil, err
	}

	// Enable it
	if err := genre.Enable(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := u.genreRepo.Update(ctx, genre); err != nil {
		return nil, err
	}

	return genre, nil
}

// DisableGenre disables a genre
func (u *genreUseCase) DisableGenre(ctx context.Context, genreID uuid.UUID) (*domain.Genre, error) {
	// Find existing genre
	genre, err := u.genreRepo.FindByID(ctx, valueobject.UUID(genreID.String()))
	if err != nil {
		return nil, err
	}

	// Disable it
	if err := genre.Disable(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := u.genreRepo.Update(ctx, genre); err != nil {
		return nil, err
	}

	return genre, nil
}