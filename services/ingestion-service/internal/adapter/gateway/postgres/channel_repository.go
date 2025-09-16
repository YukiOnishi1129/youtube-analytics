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
		ID:                id,
		YoutubeChannelID:  string(ch.YouTubeChannelID),
		Title:             ch.Title,
		ThumbnailUrl:      ch.ThumbnailURL,
		Description:       sql.NullString{String: ch.Description, Valid: ch.Description != ""},
		Country:           sql.NullString{String: ch.Country, Valid: ch.Country != ""},
		ViewCount:         sql.NullInt64{Int64: ch.ViewCount, Valid: true},
		SubscriptionCount: sql.NullInt64{Int64: ch.SubscriptionCount, Valid: true},
		VideoCount:        sql.NullInt64{Int64: ch.VideoCount, Valid: true},
		Subscribed:        sql.NullBool{Bool: ch.Subscribed, Valid: true},
		CreatedAt:         sql.NullTime{Time: ch.CreatedAt, Valid: true},
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
		ID:                id,
		Title:             ch.Title,
		ThumbnailUrl:      ch.ThumbnailURL,
		Description:       sql.NullString{String: ch.Description, Valid: ch.Description != ""},
		Country:           sql.NullString{String: ch.Country, Valid: ch.Country != ""},
		ViewCount:         sql.NullInt64{Int64: ch.ViewCount, Valid: true},
		SubscriptionCount: sql.NullInt64{Int64: ch.SubscriptionCount, Valid: true},
		VideoCount:        sql.NullInt64{Int64: ch.VideoCount, Valid: true},
		Subscribed:        sql.NullBool{Bool: ch.Subscribed, Valid: true},
		UpdatedAt:         updatedAt,
		DeletedAt:         deletedAt,
	})
}

// SaveWithSnapshots saves channel and its new snapshots in a transaction
func (r *channelRepository) SaveWithSnapshots(ctx context.Context, ch *domain.Channel) error {
	// TODO: Implement transaction handling
	// For now, just save channel
	if err := r.Save(ctx, ch); err != nil {
		return err
	}

	// Save snapshots
	for _, snapshot := range ch.GetNewSnapshots() {
		// TODO: Implement snapshot saving logic
		// This would typically be done in a transaction
		_ = snapshot
	}

	// Clear new snapshots after saving
	ch.ClearNewSnapshots()

	return nil
}

// GetByID gets a channel by ID
func (r *channelRepository) GetByID(ctx context.Context, id valueobject.UUID) (*domain.Channel, error) {
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

	return toDomainChannelFromRow(row), nil
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

	return toDomainChannelFromRow(row), nil
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

	return toDomainChannelFromYouTubeRow(row), nil
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
		channels[i] = toDomainChannelFromListRow(row)
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
		channels[i] = toDomainChannelFromListRow(row)
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
		channels[i] = toDomainChannelFromListRow(row)
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

// toDomainChannel converts database row to domain channel
func toDomainChannel(row sqlcgen.IngestionChannel) *domain.Channel {
	ch := &domain.Channel{
		ID:               valueobject.UUID(row.ID.String()),
		YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
		Title:            row.Title,
		ThumbnailURL:     row.ThumbnailUrl,
		Subscribed:       row.Subscribed.Bool,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        nullTimeToPtr(row.UpdatedAt),
	}

	if row.Description.Valid {
		ch.Description = row.Description.String
	}
	if row.Country.Valid {
		ch.Country = row.Country.String
	}
	if row.ViewCount.Valid {
		ch.ViewCount = row.ViewCount.Int64
	}
	if row.SubscriptionCount.Valid {
		ch.SubscriptionCount = row.SubscriptionCount.Int64
	}
	if row.VideoCount.Valid {
		ch.VideoCount = row.VideoCount.Int64
	}
	if row.DeletedAt.Valid {
		ch.DeletedAt = &row.DeletedAt.Time
	}

	return ch
}

// toDomainChannelFromRow converts GetChannelByIDRow to domain channel
func toDomainChannelFromRow(row sqlcgen.GetChannelByIDRow) *domain.Channel {
	ch := &domain.Channel{
		ID:               valueobject.UUID(row.ID.String()),
		YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
		Title:            row.Title,
		ThumbnailURL:     row.ThumbnailUrl,
		Subscribed:       row.Subscribed.Bool,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        nullTimeToPtr(row.UpdatedAt),
	}

	if row.Description.Valid {
		ch.Description = row.Description.String
	}
	if row.Country.Valid {
		ch.Country = row.Country.String
	}
	if row.ViewCount.Valid {
		ch.ViewCount = row.ViewCount.Int64
	}
	if row.SubscriptionCount.Valid {
		ch.SubscriptionCount = row.SubscriptionCount.Int64
	}
	if row.VideoCount.Valid {
		ch.VideoCount = row.VideoCount.Int64
	}

	return ch
}

// toDomainChannelFromYouTubeRow converts GetChannelByYouTubeIDRow to domain channel
func toDomainChannelFromYouTubeRow(row sqlcgen.GetChannelByYouTubeIDRow) *domain.Channel {
	ch := &domain.Channel{
		ID:               valueobject.UUID(row.ID.String()),
		YouTubeChannelID: valueobject.YouTubeChannelID(row.YoutubeChannelID),
		Title:            row.Title,
		ThumbnailURL:     row.ThumbnailUrl,
		Subscribed:       row.Subscribed.Bool,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        nullTimeToPtr(row.UpdatedAt),
	}

	if row.Description.Valid {
		ch.Description = row.Description.String
	}
	if row.Country.Valid {
		ch.Country = row.Country.String
	}
	if row.ViewCount.Valid {
		ch.ViewCount = row.ViewCount.Int64
	}
	if row.SubscriptionCount.Valid {
		ch.SubscriptionCount = row.SubscriptionCount.Int64
	}
	if row.VideoCount.Valid {
		ch.VideoCount = row.VideoCount.Int64
	}

	return ch
}

// toDomainChannelFromListRow converts ListChannelsRow, ListActiveChannelsRow, ListSubscribedChannelsRow to domain channel
func toDomainChannelFromListRow(row interface{}) *domain.Channel {
	switch r := row.(type) {
	case sqlcgen.ListChannelsRow:
		return &domain.Channel{
			ID:               valueobject.UUID(r.ID.String()),
			YouTubeChannelID: valueobject.YouTubeChannelID(r.YoutubeChannelID),
			Title:            r.Title,
			ThumbnailURL:     r.ThumbnailUrl,
			Subscribed:       r.Subscribed.Bool,
			CreatedAt:        r.CreatedAt.Time,
			UpdatedAt:        nullTimeToPtr(r.UpdatedAt),
			Description:      nullStringToString(r.Description),
			Country:          nullStringToString(r.Country),
			ViewCount:        nullInt64ToInt64(r.ViewCount),
			SubscriptionCount: nullInt64ToInt64(r.SubscriptionCount),
			VideoCount:        nullInt64ToInt64(r.VideoCount),
		}
	case sqlcgen.ListActiveChannelsRow:
		return &domain.Channel{
			ID:               valueobject.UUID(r.ID.String()),
			YouTubeChannelID: valueobject.YouTubeChannelID(r.YoutubeChannelID),
			Title:            r.Title,
			ThumbnailURL:     r.ThumbnailUrl,
			Subscribed:       r.Subscribed.Bool,
			CreatedAt:        r.CreatedAt.Time,
			UpdatedAt:        nullTimeToPtr(r.UpdatedAt),
			Description:      nullStringToString(r.Description),
			Country:          nullStringToString(r.Country),
			ViewCount:        nullInt64ToInt64(r.ViewCount),
			SubscriptionCount: nullInt64ToInt64(r.SubscriptionCount),
			VideoCount:        nullInt64ToInt64(r.VideoCount),
		}
	case sqlcgen.ListSubscribedChannelsRow:
		return &domain.Channel{
			ID:               valueobject.UUID(r.ID.String()),
			YouTubeChannelID: valueobject.YouTubeChannelID(r.YoutubeChannelID),
			Title:            r.Title,
			ThumbnailURL:     r.ThumbnailUrl,
			Subscribed:       r.Subscribed.Bool,
			CreatedAt:        r.CreatedAt.Time,
			UpdatedAt:        nullTimeToPtr(r.UpdatedAt),
			Description:      nullStringToString(r.Description),
			Country:          nullStringToString(r.Country),
			ViewCount:        nullInt64ToInt64(r.ViewCount),
			SubscriptionCount: nullInt64ToInt64(r.SubscriptionCount),
			VideoCount:        nullInt64ToInt64(r.VideoCount),
		}
	default:
		panic("unsupported row type")
	}
}

// nullStringToString converts sql.NullString to string
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// nullInt64ToInt64 converts sql.NullInt64 to int64
func nullInt64ToInt64(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}