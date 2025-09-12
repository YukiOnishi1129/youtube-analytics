package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/postgres/sqlcgen"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"github.com/google/uuid"
)

// channelRepository implements gateway.ChannelRepository interface
type channelRepository struct {
	*Repository
}

// NewChannelRepository creates a new channel repository
func NewChannelRepository(repo *Repository) gateway.ChannelRepository {
	return &channelRepository{Repository: repo}
}

// Save creates a new channel
func (r *channelRepository) Save(ctx context.Context, ch *domain.Channel) error {
	id, err := uuid.Parse(string(ch.ID))
	if err != nil {
		return err
	}

	return r.q.CreateChannel(ctx, sqlcgen.CreateChannelParams{
		ID:               id,
		YoutubeChannelID: string(ch.YouTubeChannelID),
		Title:            ch.Title,
		ThumbnailUrl:     ch.ThumbnailURL,
		Subscribed:       sql.NullBool{Bool: ch.Subscribed, Valid: true},
		CreatedAt:        sql.NullTime{Time: ch.CreatedAt, Valid: true},
	})
}

// Update updates an existing channel
func (r *channelRepository) Update(ctx context.Context, ch *domain.Channel) error {
	id, err := uuid.Parse(string(ch.ID))
	if err != nil {
		return err
	}

	var deletedAt sql.NullTime
	if ch.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *ch.DeletedAt, Valid: true}
	}

	var updatedAt sql.NullTime
	if ch.UpdatedAt != nil {
		updatedAt = sql.NullTime{Time: *ch.UpdatedAt, Valid: true}
	}

	return r.q.UpdateChannel(ctx, sqlcgen.UpdateChannelParams{
		ID:           id,
		Title:        ch.Title,
		ThumbnailUrl: ch.ThumbnailURL,
		Subscribed:   sql.NullBool{Bool: ch.Subscribed, Valid: true},
		UpdatedAt:    updatedAt,
		DeletedAt:    deletedAt,
	})
}

// FindByID finds a channel by ID
func (r *channelRepository) FindByID(ctx context.Context, id valueobject.UUID) (*domain.Channel, error) {
	uid, err := uuid.Parse(string(id))
	if err != nil {
		return nil, err
	}

	row, err := r.q.GetChannelByID(ctx, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrChannelNotFound
		}
		return nil, err
	}

	return &domain.Channel{
		ID:               valueobject.UUID(row.ID.String()),
		YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
		Title:            row.Title,
		ThumbnailURL:     row.ThumbnailUrl,
		Subscribed:       row.Subscribed.Bool,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        nullTimeToPtr(row.UpdatedAt),
	}, nil
}

// FindByYouTubeID finds a channel by YouTube ID
func (r *channelRepository) FindByYouTubeID(ctx context.Context, ytID valueobject.YouTubeChannelID) (*domain.Channel, error) {
	row, err := r.q.GetChannelByYouTubeID(ctx, string(ytID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrChannelNotFound
		}
		return nil, err
	}

	return &domain.Channel{
		ID:               valueobject.UUID(row.ID.String()),
		YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
		Title:            row.Title,
		ThumbnailURL:     row.ThumbnailUrl,
		Subscribed:       row.Subscribed.Bool,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        nullTimeToPtr(row.UpdatedAt),
	}, nil
}

// FindByYouTubeChannelID finds a channel by YouTube channel ID
func (r *channelRepository) FindByYouTubeChannelID(ctx context.Context, youtubeChannelID valueobject.YouTubeChannelID) (*domain.Channel, error) {
	return r.FindByYouTubeID(ctx, youtubeChannelID)
}

// List lists channels with pagination
func (r *channelRepository) List(ctx context.Context, subscribed *bool, q *string, sort string, limit, offset int) ([]*domain.Channel, error) {
	// TODO: Implement filtering by subscribed and query
	// For now, just return all channels with pagination
	rows, err := r.q.ListChannels(ctx, sqlcgen.ListChannelsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	channels := make([]*domain.Channel, len(rows))
	for i, row := range rows {
		channels[i] = &domain.Channel{
			ID:               valueobject.UUID(row.ID.String()),
			YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
			Title:            row.Title,
			ThumbnailURL:     row.ThumbnailUrl,
			Subscribed:       row.Subscribed.Bool,
			CreatedAt:        row.CreatedAt.Time,
			UpdatedAt:        nullTimeToPtr(row.UpdatedAt),
		}
	}
	return channels, nil
}

// Count counts channels
func (r *channelRepository) Count(ctx context.Context, subscribed *bool, q *string) (int, error) {
	// TODO: Implement filtering by subscribed and query
	count, err := r.q.CountChannels(ctx)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// ListActive lists all active channels
func (r *channelRepository) ListActive(ctx context.Context) ([]*domain.Channel, error) {
	rows, err := r.q.ListActiveChannels(ctx)
	if err != nil {
		return nil, err
	}

	channels := make([]*domain.Channel, len(rows))
	for i, row := range rows {
		channels[i] = &domain.Channel{
			ID:               valueobject.UUID(row.ID.String()),
			YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
			Title:            row.Title,
			ThumbnailURL:     row.ThumbnailUrl,
			Subscribed:       row.Subscribed.Bool,
			CreatedAt:        row.CreatedAt.Time,
			UpdatedAt:        nullTimeToPtr(row.UpdatedAt),
		}
	}
	return channels, nil
}

// ListSubscribed lists all subscribed channels
func (r *channelRepository) ListSubscribed(ctx context.Context) ([]*domain.Channel, error) {
	rows, err := r.q.ListSubscribedChannels(ctx)
	if err != nil {
		return nil, err
	}

	channels := make([]*domain.Channel, len(rows))
	for i, row := range rows {
		channels[i] = &domain.Channel{
			ID:               valueobject.UUID(row.ID.String()),
			YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
			Title:            row.Title,
			ThumbnailURL:     row.ThumbnailUrl,
			Subscribed:       row.Subscribed.Bool,
			CreatedAt:        row.CreatedAt.Time,
			UpdatedAt:        nullTimeToPtr(row.UpdatedAt),
		}
	}
	return channels, nil
}

// nullTimeToPtr converts sql.NullTime to *time.Time
func nullTimeToPtr(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
