package gateway

import (
	"context"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// KeywordRepository is the repository interface for Keyword aggregate
type KeywordRepository interface {
	Save(ctx context.Context, k *domain.Keyword) error
	FindAll(ctx context.Context, enabledOnly bool) ([]*domain.Keyword, error)
	FindByID(ctx context.Context, id valueobject.UUID) (*domain.Keyword, error)
	SoftDelete(ctx context.Context, id valueobject.UUID) error
}

// ChannelRepository is the repository interface for Channel aggregate
type ChannelRepository interface {
	Save(ctx context.Context, ch *domain.Channel) error
	Update(ctx context.Context, ch *domain.Channel) error
	FindByID(ctx context.Context, id valueobject.UUID) (*domain.Channel, error)
	FindByYouTubeID(ctx context.Context, ytID valueobject.YouTubeChannelID) (*domain.Channel, error)
	FindByYouTubeChannelID(ctx context.Context, youtubeChannelID valueobject.YouTubeChannelID) (*domain.Channel, error)
	List(ctx context.Context, subscribed *bool, q *string, sort string, limit, offset int) ([]*domain.Channel, error)
	Count(ctx context.Context, subscribed *bool, q *string) (int, error)
	ListActive(ctx context.Context) ([]*domain.Channel, error)
	ListSubscribed(ctx context.Context) ([]*domain.Channel, error)
}

// ChannelSnapshotRepository is the repository interface for ChannelSnapshot
type ChannelSnapshotRepository interface {
	Append(ctx context.Context, snap *domain.ChannelSnapshot) error
	Latest(ctx context.Context, channelID valueobject.UUID) (*domain.ChannelSnapshot, error)
	ListByChannel(ctx context.Context, channelID valueobject.UUID, limit int) ([]*domain.ChannelSnapshot, error)
}

// VideoRepository is the repository interface for Video aggregate
type VideoRepository interface {
	Save(ctx context.Context, v *domain.Video) error
	FindByID(ctx context.Context, id valueobject.UUID) (*domain.Video, error)
	FindByYouTubeID(ctx context.Context, ytID valueobject.YouTubeVideoID) (*domain.Video, error)
	ExistsByYouTubeVideoID(ctx context.Context, youtubeVideoID valueobject.YouTubeVideoID) (bool, error)
	ListByChannel(ctx context.Context, channelID valueobject.UUID, limit, offset int) ([]*domain.Video, error)
	CountByChannel(ctx context.Context, channelID valueobject.UUID) (int, error)
	ListActive(ctx context.Context, since time.Time) ([]*domain.Video, error)
}

// VideoSnapshotRepository is the repository interface for VideoSnapshot
type VideoSnapshotRepository interface {
	Insert(ctx context.Context, s *domain.VideoSnapshot) error
	Exists(ctx context.Context, videoID valueobject.UUID, cp valueobject.CheckpointHour) (bool, error)
	FindByVideoAndCP(ctx context.Context, videoID valueobject.UUID, cp valueobject.CheckpointHour) (*domain.VideoSnapshot, error)
	ListByVideo(ctx context.Context, videoID valueobject.UUID) ([]*domain.VideoSnapshot, error)
}