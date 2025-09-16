package postgres

import (
	"context"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/postgres/sqlcgen"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"github.com/google/uuid"
)

// videoGenreRepository implements gateway.VideoGenreRepository interface
type videoGenreRepository struct {
	*Repository
}

// NewVideoGenreRepository creates a new video genre repository
func NewVideoGenreRepository(repo *Repository) gateway.VideoGenreRepository {
	return &videoGenreRepository{Repository: repo}
}

// Save creates a new video-genre association
func (r *videoGenreRepository) Save(ctx context.Context, vg *domain.VideoGenre) error {
	id, err := uuid.Parse(string(vg.ID))
	if err != nil {
		return err
	}

	videoID, err := uuid.Parse(string(vg.VideoID))
	if err != nil {
		return err
	}

	genreID, err := uuid.Parse(string(vg.GenreID))
	if err != nil {
		return err
	}

	return r.q.CreateVideoGenre(ctx, sqlcgen.CreateVideoGenreParams{
		ID:        id,
		VideoID:   videoID,
		GenreID:   genreID,
		CreatedAt: time.Now(),
	})
}

// SaveBatch creates multiple video-genre associations
func (r *videoGenreRepository) SaveBatch(ctx context.Context, vgs []*domain.VideoGenre) error {
	return r.ExecTx(ctx, func(repo *Repository) error {
		for _, vg := range vgs {
			id, err := uuid.Parse(string(vg.ID))
			if err != nil {
				return err
			}

			videoID, err := uuid.Parse(string(vg.VideoID))
			if err != nil {
				return err
			}

			genreID, err := uuid.Parse(string(vg.GenreID))
			if err != nil {
				return err
			}

			if err := repo.q.CreateVideoGenre(ctx, sqlcgen.CreateVideoGenreParams{
				ID:        id,
				VideoID:   videoID,
				GenreID:   genreID,
				CreatedAt: vg.CreatedAt,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

// FindByVideo finds all video-genre associations for a video
func (r *videoGenreRepository) FindByVideo(ctx context.Context, videoID valueobject.UUID) ([]*domain.VideoGenre, error) {
	vid, err := uuid.Parse(string(videoID))
	if err != nil {
		return nil, err
	}

	rows, err := r.q.ListVideoGenresByVideo(ctx, vid)
	if err != nil {
		return nil, err
	}

	videoGenres := make([]*domain.VideoGenre, len(rows))
	for i, row := range rows {
		videoGenres[i] = toDomainVideoGenre(row)
	}

	return videoGenres, nil
}

// FindByGenre finds all video-genre associations for a genre
func (r *videoGenreRepository) FindByGenre(ctx context.Context, genreID valueobject.UUID) ([]*domain.VideoGenre, error) {
	gid, err := uuid.Parse(string(genreID))
	if err != nil {
		return nil, err
	}

	rows, err := r.q.ListVideoGenresByGenre(ctx, gid)
	if err != nil {
		return nil, err
	}

	videoGenres := make([]*domain.VideoGenre, len(rows))
	for i, row := range rows {
		videoGenres[i] = toDomainVideoGenre(row)
	}

	return videoGenres, nil
}

// ExistsByVideoAndGenre checks if a video-genre association exists
func (r *videoGenreRepository) ExistsByVideoAndGenre(ctx context.Context, videoID, genreID valueobject.UUID) (bool, error) {
	vid, err := uuid.Parse(string(videoID))
	if err != nil {
		return false, err
	}

	gid, err := uuid.Parse(string(genreID))
	if err != nil {
		return false, err
	}

	return r.q.CheckVideoGenreExists(ctx, sqlcgen.CheckVideoGenreExistsParams{
		VideoID: vid,
		GenreID: gid,
	})
}

// DeleteByVideo deletes all video-genre associations for a video
func (r *videoGenreRepository) DeleteByVideo(ctx context.Context, videoID valueobject.UUID) error {
	vid, err := uuid.Parse(string(videoID))
	if err != nil {
		return err
	}

	return r.q.DeleteVideoGenresByVideo(ctx, vid)
}

// DeleteByGenre deletes all video-genre associations for a genre
func (r *videoGenreRepository) DeleteByGenre(ctx context.Context, genreID valueobject.UUID) error {
	gid, err := uuid.Parse(string(genreID))
	if err != nil {
		return err
	}

	return r.q.DeleteVideoGenresByGenre(ctx, gid)
}

// toDomainVideoGenre converts a database row to a domain video genre
func toDomainVideoGenre(row sqlcgen.IngestionVideoGenre) *domain.VideoGenre {
	return &domain.VideoGenre{
		ID:        valueobject.UUID(row.ID.String()),
		VideoID:   valueobject.UUID(row.VideoID.String()),
		GenreID:   valueobject.UUID(row.GenreID.String()),
		CreatedAt: row.CreatedAt,
	}
}