package usecase

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"github.com/google/uuid"
)

// videoGenreUseCase implements the VideoGenreInputPort interface
type videoGenreUseCase struct {
	videoGenreRepo gateway.VideoGenreRepository
}

// NewVideoGenreUseCase creates a new video genre use case
func NewVideoGenreUseCase(videoGenreRepo gateway.VideoGenreRepository) input.VideoGenreInputPort {
	return &videoGenreUseCase{
		videoGenreRepo: videoGenreRepo,
	}
}

// AssociateVideoWithGenre associates a video with a genre
func (u *videoGenreUseCase) AssociateVideoWithGenre(ctx context.Context, videoID, genreID uuid.UUID) (*domain.VideoGenre, error) {
	// Check if association already exists
	exists, err := u.videoGenreRepo.ExistsByVideoAndGenre(ctx, valueobject.UUID(videoID.String()), valueobject.UUID(genreID.String()))
	if err != nil {
		return nil, err
	}
	if exists {
		// Return existing association
		associations, err := u.videoGenreRepo.FindByVideo(ctx, valueobject.UUID(videoID.String()))
		if err != nil {
			return nil, err
		}
		for _, assoc := range associations {
			if assoc.GenreID == valueobject.UUID(genreID.String()) {
				return assoc, nil
			}
		}
	}

	// Create new association
	videoGenre := &domain.VideoGenre{
		ID:      valueobject.UUID(uuid.New().String()),
		VideoID: valueobject.UUID(videoID.String()),
		GenreID: valueobject.UUID(genreID.String()),
	}

	// Save to repository
	if err := u.videoGenreRepo.Save(ctx, videoGenre); err != nil {
		return nil, err
	}

	return videoGenre, nil
}

// AssociateVideoWithGenres associates a video with multiple genres
func (u *videoGenreUseCase) AssociateVideoWithGenres(ctx context.Context, videoID uuid.UUID, genreIDs []uuid.UUID) ([]*domain.VideoGenre, error) {
	// Create associations
	videoGenres := make([]*domain.VideoGenre, 0, len(genreIDs))
	for _, genreID := range genreIDs {
		// Check if association already exists
		exists, err := u.videoGenreRepo.ExistsByVideoAndGenre(ctx, valueobject.UUID(videoID.String()), valueobject.UUID(genreID.String()))
		if err != nil {
			return nil, err
		}
		if !exists {
			videoGenre := &domain.VideoGenre{
				ID:      valueobject.UUID(uuid.New().String()),
				VideoID: valueobject.UUID(videoID.String()),
				GenreID: valueobject.UUID(genreID.String()),
			}
			videoGenres = append(videoGenres, videoGenre)
		}
	}

	// Save batch to repository
	if len(videoGenres) > 0 {
		if err := u.videoGenreRepo.SaveBatch(ctx, videoGenres); err != nil {
			return nil, err
		}
	}

	// Return all associations
	return u.videoGenreRepo.FindByVideo(ctx, valueobject.UUID(videoID.String()))
}

// GetVideoGenres gets all genres associated with a video
func (u *videoGenreUseCase) GetVideoGenres(ctx context.Context, videoID uuid.UUID) ([]*domain.VideoGenre, error) {
	return u.videoGenreRepo.FindByVideo(ctx, valueobject.UUID(videoID.String()))
}

// GetGenreVideos gets all videos associated with a genre
func (u *videoGenreUseCase) GetGenreVideos(ctx context.Context, genreID uuid.UUID) ([]*domain.VideoGenre, error) {
	return u.videoGenreRepo.FindByGenre(ctx, valueobject.UUID(genreID.String()))
}

// DisassociateVideoFromGenre removes the association between a video and a genre
func (u *videoGenreUseCase) DisassociateVideoFromGenre(ctx context.Context, videoID, genreID uuid.UUID) error {
	// For now, we don't have a specific delete method, so we'll use the general delete by video
	// In a real implementation, you might want to add a specific method to delete a single association
	return domain.ErrNotFound // Placeholder - implement specific deletion if needed
}

// DisassociateVideoFromAllGenres removes all genre associations for a video
func (u *videoGenreUseCase) DisassociateVideoFromAllGenres(ctx context.Context, videoID uuid.UUID) error {
	return u.videoGenreRepo.DeleteByVideo(ctx, valueobject.UUID(videoID.String()))
}