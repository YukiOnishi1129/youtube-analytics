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

	channelID, err := uuid.Parse(string(v.ChannelID))
	if err != nil {
		return err
	}

	return r.q.CreateVideo(ctx, sqlcgen.CreateVideoParams{
		ID:               id,
		YoutubeVideoID:   string(v.YouTubeVideoID),
		ChannelID:        channelID,
		YoutubeChannelID: string(v.YouTubeChannelID),
		Title:            v.Title,
		PublishedAt:      v.PublishedAt,
		CategoryID:       int32(v.CategoryID),
		CreatedAt:        sql.NullTime{Time: v.CreatedAt, Valid: true},
	})
}

// SaveWithSnapshots saves video and its new snapshots in a transaction
func (r *videoRepository) SaveWithSnapshots(ctx context.Context, v *domain.Video) error {
	// TODO: Implement transaction handling
	// For now, just save video
	if err := r.Save(ctx, v); err != nil {
		return err
	}

	// Save snapshots
	for _, snapshot := range v.GetNewSnapshots() {
		// TODO: Implement snapshot saving logic
		// This would typically be done in a transaction
		_ = snapshot
	}

	// Clear new snapshots after saving
	v.ClearNewSnapshots()
	
	return nil
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

	return toDomainVideoFromRow(row), nil
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

	return toDomainVideoFromRow(row), nil
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

	return toDomainVideoFromYouTubeRow(row), nil
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
		videos[i] = toDomainVideo(row.ID, row.ChannelID, row.YoutubeVideoID, row.YoutubeChannelID, row.Title,
			row.PublishedAt, row.CategoryID, row.CreatedAt, sql.NullTime{}, sql.NullTime{})
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
		videos[i] = toDomainVideo(row.ID, row.ChannelID, row.YoutubeVideoID, row.YoutubeChannelID, row.Title,
			row.PublishedAt, row.CategoryID, row.CreatedAt, sql.NullTime{}, sql.NullTime{})
	}
	return videos, nil
}

// toDomainVideo converts database row to domain video
func toDomainVideo(id uuid.UUID, channelID uuid.UUID, youtubeVideoID, youtubeChannelID, title string, 
	publishedAt time.Time, categoryID int32, createdAt sql.NullTime, updatedAt, deletedAt sql.NullTime) *domain.Video {
	v := &domain.Video{
		ID:               valueobject.UUID(id.String()),
		ChannelID:        valueobject.UUID(channelID.String()),
		YouTubeVideoID:   valueobject.YouTubeVideoID(youtubeVideoID),
		YouTubeChannelID: valueobject.YouTubeChannelID(youtubeChannelID),
		Title:            title,
		PublishedAt:      publishedAt,
		CategoryID:       valueobject.CategoryID(categoryID),
		CreatedAt:        createdAt.Time,
	}

	if updatedAt.Valid {
		v.UpdatedAt = &updatedAt.Time
	}
	if deletedAt.Valid {
		v.DeletedAt = &deletedAt.Time
	}

	return v
}

// toDomainVideoFromRow converts GetVideoByIDRow to domain video
func toDomainVideoFromRow(row sqlcgen.GetVideoByIDRow) *domain.Video {
	return &domain.Video{
		ID:               valueobject.UUID(row.ID.String()),
		ChannelID:        valueobject.UUID(row.ChannelID.String()),
		YouTubeVideoID:   valueobject.YouTubeVideoID(row.YoutubeVideoID),
		YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
		Title:            row.Title,
		PublishedAt:      row.PublishedAt,
		CategoryID:       valueobject.CategoryID(row.CategoryID),
		CreatedAt:        row.CreatedAt.Time,
	}
}

// toDomainVideoFromYouTubeRow converts GetVideoByYouTubeIDRow to domain video  
func toDomainVideoFromYouTubeRow(row sqlcgen.GetVideoByYouTubeIDRow) *domain.Video {
	return &domain.Video{
		ID:               valueobject.UUID(row.ID.String()),
		ChannelID:        valueobject.UUID(row.ChannelID.String()),
		YouTubeVideoID:   valueobject.YouTubeVideoID(row.YoutubeVideoID),
		YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
		Title:            row.Title,
		PublishedAt:      row.PublishedAt,
		CategoryID:       valueobject.CategoryID(row.CategoryID),
		CreatedAt:        row.CreatedAt.Time,
	}
}