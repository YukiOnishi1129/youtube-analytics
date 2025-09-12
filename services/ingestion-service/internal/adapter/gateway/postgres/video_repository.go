package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/postgres/sqlcgen"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

// videoRepository implements gateway.VideoRepository interface
type videoRepository struct {
	*Repository
}

// NewVideoRepository creates a new video repository
func NewVideoRepository(repo *Repository) gateway.VideoRepository {
	return &videoRepository{Repository: repo}
}

// Save creates a new video
func (r *videoRepository) Save(ctx context.Context, v *domain.Video) error {
	id, err := uuid.Parse(string(v.ID))
	if err != nil {
		return err
	}

	// TODO: Get channel_id from channels table based on youtube_channel_id
	// For now, using a placeholder
	channelID := uuid.New()
	videoURL := "https://www.youtube.com/watch?v=" + string(v.YouTubeVideoID)

	return r.q.CreateVideo(ctx, sqlcgen.CreateVideoParams{
		ID:               id,
		YoutubeVideoID:   string(v.YouTubeVideoID),
		ChannelID:        channelID,
		YoutubeChannelID: string(v.YouTubeChannelID),
		Title:            v.Title,
		PublishedAt:      v.PublishedAt,
		CategoryID:       0, // Default category
		ThumbnailUrl:     v.ThumbnailURL,
		VideoUrl:         videoURL,
		CreatedAt:        sql.NullTime{Time: v.CreatedAt, Valid: true},
	})
}

// GetByID gets a video by ID
func (r *videoRepository) GetByID(ctx context.Context, id valueobject.UUID) (*domain.Video, error) {
	uid, err := uuid.Parse(string(id))
	if err != nil {
		return nil, err
	}

	row, err := r.q.GetVideoByID(ctx, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrVideoNotFound
		}
		return nil, err
	}

	return &domain.Video{
		ID:               valueobject.UUID(row.ID.String()),
		YouTubeVideoID:   valueobject.YouTubeVideoID(row.YoutubeVideoID),
		YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
		Title:            row.Title,
		ThumbnailURL:     row.ThumbnailUrl,
		PublishedAt:      row.PublishedAt,
		CreatedAt:        row.CreatedAt.Time,
	}, nil
}

// FindByID finds a video by ID
func (r *videoRepository) FindByID(ctx context.Context, id valueobject.UUID) (*domain.Video, error) {
	uid, err := uuid.Parse(string(id))
	if err != nil {
		return nil, err
	}

	row, err := r.q.GetVideoByID(ctx, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrVideoNotFound
		}
		return nil, err
	}

	return &domain.Video{
		ID:               valueobject.UUID(row.ID.String()),
		YouTubeVideoID:   valueobject.YouTubeVideoID(row.YoutubeVideoID),
		YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
		Title:            row.Title,
		ThumbnailURL:     row.ThumbnailUrl,
		PublishedAt:      row.PublishedAt,
		CreatedAt:        row.CreatedAt.Time,
	}, nil
}

// FindByYouTubeID finds a video by YouTube ID
func (r *videoRepository) FindByYouTubeID(ctx context.Context, ytID valueobject.YouTubeVideoID) (*domain.Video, error) {
	row, err := r.q.GetVideoByYouTubeID(ctx, string(ytID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrVideoNotFound
		}
		return nil, err
	}

	return &domain.Video{
		ID:               valueobject.UUID(row.ID.String()),
		YouTubeVideoID:   valueobject.YouTubeVideoID(row.YoutubeVideoID),
		YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
		Title:            row.Title,
		ThumbnailURL:     row.ThumbnailUrl,
		PublishedAt:      row.PublishedAt,
		CreatedAt:        row.CreatedAt.Time,
	}, nil
}

// ExistsByYouTubeVideoID checks if a video exists by YouTube ID
func (r *videoRepository) ExistsByYouTubeVideoID(ctx context.Context, youtubeVideoID valueobject.YouTubeVideoID) (bool, error) {
	return r.q.CheckVideoExists(ctx, string(youtubeVideoID))
}

// ListByChannel lists videos by channel with pagination
func (r *videoRepository) ListByChannel(ctx context.Context, channelID valueobject.UUID, limit, offset int) ([]*domain.Video, error) {
	uid, err := uuid.Parse(string(channelID))
	if err != nil {
		return nil, err
	}

	rows, err := r.q.ListVideosByChannel(ctx, sqlcgen.ListVideosByChannelParams{
		ChannelID: uid,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return nil, err
	}

	videos := make([]*domain.Video, len(rows))
	for i, row := range rows {
		videos[i] = &domain.Video{
			ID:               valueobject.UUID(row.ID.String()),
			YouTubeVideoID:   valueobject.YouTubeVideoID(row.YoutubeVideoID),
			YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
			Title:            row.Title,
			ThumbnailURL:     row.ThumbnailUrl,
			PublishedAt:      row.PublishedAt,
			CreatedAt:        row.CreatedAt.Time,
		}
	}
	return videos, nil
}

// CountByChannel counts videos in a channel
func (r *videoRepository) CountByChannel(ctx context.Context, channelID valueobject.UUID) (int, error) {
	uid, err := uuid.Parse(string(channelID))
	if err != nil {
		return 0, err
	}

	count, err := r.q.CountVideosByChannel(ctx, uid)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// ListActive lists active videos published after the given time
func (r *videoRepository) ListActive(ctx context.Context, since time.Time) ([]*domain.Video, error) {
	rows, err := r.q.ListActiveVideos(ctx, since)
	if err != nil {
		return nil, err
	}

	videos := make([]*domain.Video, len(rows))
	for i, row := range rows {
		videos[i] = &domain.Video{
			ID:               valueobject.UUID(row.ID.String()),
			YouTubeVideoID:   valueobject.YouTubeVideoID(row.YoutubeVideoID),
			YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
			Title:            row.Title,
			ThumbnailURL:     row.ThumbnailUrl,
			PublishedAt:      row.PublishedAt,
			CreatedAt:        row.CreatedAt.Time,
		}
	}
	return videos, nil
}