package input

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/google/uuid"
)

// GenreInputPort is the interface for genre use cases
type GenreInputPort interface {
	ListGenres(ctx context.Context) ([]*domain.Genre, error)
	ListEnabledGenres(ctx context.Context) ([]*domain.Genre, error)
	GetGenre(ctx context.Context, genreID uuid.UUID) (*domain.Genre, error)
	GetGenreByCode(ctx context.Context, code string) (*domain.Genre, error)
	CreateGenre(ctx context.Context, input *CreateGenreInput) (*domain.Genre, error)
	UpdateGenre(ctx context.Context, input *UpdateGenreInput) (*domain.Genre, error)
	EnableGenre(ctx context.Context, genreID uuid.UUID) (*domain.Genre, error)
	DisableGenre(ctx context.Context, genreID uuid.UUID) (*domain.Genre, error)
}

// CreateGenreInput represents the input for creating a genre
type CreateGenreInput struct {
	Code        string
	Name        string
	Language    string
	RegionCode  string
	CategoryIDs []int
}

// UpdateGenreInput represents the input for updating a genre
type UpdateGenreInput struct {
	GenreID     uuid.UUID
	Name        string
	CategoryIDs []int
}