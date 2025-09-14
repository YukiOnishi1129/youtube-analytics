package input

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/google/uuid"
)

// VideoGenreInputPort is the interface for video-genre association use cases
type VideoGenreInputPort interface {
	AssociateVideoWithGenre(ctx context.Context, videoID, genreID uuid.UUID) (*domain.VideoGenre, error)
	AssociateVideoWithGenres(ctx context.Context, videoID uuid.UUID, genreIDs []uuid.UUID) ([]*domain.VideoGenre, error)
	GetVideoGenres(ctx context.Context, videoID uuid.UUID) ([]*domain.VideoGenre, error)
	GetGenreVideos(ctx context.Context, genreID uuid.UUID) ([]*domain.VideoGenre, error)
	DisassociateVideoFromGenre(ctx context.Context, videoID, genreID uuid.UUID) error
	DisassociateVideoFromAllGenres(ctx context.Context, videoID uuid.UUID) error
}