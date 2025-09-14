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
	Update(ctx context.Context, k *domain.Keyword) error
	FindAll(ctx context.Context, enabledOnly bool) ([]*domain.Keyword, error)
	FindByID(ctx context.Context, id valueobject.UUID) (*domain.Keyword, error)
	FindByGenre(ctx context.Context, genreID valueobject.UUID, enabledOnly bool) ([]*domain.Keyword, error)
	FindByGenreAndType(ctx context.Context, genreID valueobject.UUID, filterType valueobject.FilterType, enabledOnly bool) ([]*domain.Keyword, error)
	SoftDelete(ctx context.Context, id valueobject.UUID) error
}

// ChannelRepository is the repository interface for Channel aggregate
type ChannelRepository interface {
	Save(ctx context.Context, ch *domain.Channel) error
	SaveWithSnapshots(ctx context.Context, ch *domain.Channel) error // Save channel and its new snapshots
	Update(ctx context.Context, ch *domain.Channel) error
	GetByID(ctx context.Context, id valueobject.UUID) (*domain.Channel, error)
	FindByID(ctx context.Context, id valueobject.UUID) (*domain.Channel, error)
	FindByYouTubeID(ctx context.Context, ytID valueobject.YouTubeChannelID) (*domain.Channel, error)
	FindByYouTubeChannelID(ctx context.Context, youtubeChannelID valueobject.YouTubeChannelID) (*domain.Channel, error)
	List(ctx context.Context, subscribed *bool, q *string, sort string, limit, offset int) ([]*domain.Channel, error)
	Count(ctx context.Context, subscribed *bool, q *string) (int, error)
	ListActive(ctx context.Context) ([]*domain.Channel, error)
	ListSubscribed(ctx context.Context) ([]*domain.Channel, error)
}

// ChannelSnapshotRepository is the repository interface for ChannelSnapshot (read-only)
type ChannelSnapshotRepository interface {
	Latest(ctx context.Context, channelID valueobject.UUID) (*domain.ChannelSnapshot, error)
	ListByChannel(ctx context.Context, channelID valueobject.UUID, limit int) ([]*domain.ChannelSnapshot, error)
}

// VideoRepository is the repository interface for Video aggregate
type VideoRepository interface {
	Save(ctx context.Context, v *domain.Video) error
	SaveWithSnapshots(ctx context.Context, v *domain.Video) error // Save video and its new snapshots
	GetByID(ctx context.Context, id valueobject.UUID) (*domain.Video, error)
	FindByID(ctx context.Context, id valueobject.UUID) (*domain.Video, error)
	FindByYouTubeID(ctx context.Context, ytID valueobject.YouTubeVideoID) (*domain.Video, error)
	ExistsByYouTubeVideoID(ctx context.Context, youtubeVideoID valueobject.YouTubeVideoID) (bool, error)
	ListByChannel(ctx context.Context, channelID valueobject.UUID, limit, offset int) ([]*domain.Video, error)
	CountByChannel(ctx context.Context, channelID valueobject.UUID) (int, error)
	ListActive(ctx context.Context, since time.Time) ([]*domain.Video, error)
}

// VideoSnapshotRepository is the repository interface for VideoSnapshot (read-only)
type VideoSnapshotRepository interface {
	Exists(ctx context.Context, videoID valueobject.UUID, cp valueobject.CheckpointHour) (bool, error)
	FindByVideoAndCP(ctx context.Context, videoID valueobject.UUID, cp valueobject.CheckpointHour) (*domain.VideoSnapshot, error)
	ListByVideo(ctx context.Context, videoID valueobject.UUID) ([]*domain.VideoSnapshot, error)
	ListByVideoID(ctx context.Context, videoID valueobject.UUID) ([]*domain.VideoSnapshot, error)
}

// GenreRepository is the repository interface for Genre aggregate
type GenreRepository interface {
	Save(ctx context.Context, g *domain.Genre) error
	Update(ctx context.Context, g *domain.Genre) error
	FindByID(ctx context.Context, id valueobject.UUID) (*domain.Genre, error)
	FindByCode(ctx context.Context, code string) (*domain.Genre, error)
	FindAll(ctx context.Context) ([]*domain.Genre, error)
	FindEnabled(ctx context.Context) ([]*domain.Genre, error)
}

// YouTubeCategoryRepository is the repository interface for YouTubeCategory aggregate
type YouTubeCategoryRepository interface {
	Save(ctx context.Context, c *domain.YouTubeCategory) error
	Update(ctx context.Context, c *domain.YouTubeCategory) error
	FindByID(ctx context.Context, id valueobject.CategoryID) (*domain.YouTubeCategory, error)
	FindAll(ctx context.Context) ([]*domain.YouTubeCategory, error)
	FindAssignable(ctx context.Context) ([]*domain.YouTubeCategory, error)
}

// VideoGenreRepository is the repository interface for VideoGenre relationships
type VideoGenreRepository interface {
	Save(ctx context.Context, vg *domain.VideoGenre) error
	SaveBatch(ctx context.Context, vgs []*domain.VideoGenre) error
	FindByVideo(ctx context.Context, videoID valueobject.UUID) ([]*domain.VideoGenre, error)
	FindByGenre(ctx context.Context, genreID valueobject.UUID) ([]*domain.VideoGenre, error)
	ExistsByVideoAndGenre(ctx context.Context, videoID, genreID valueobject.UUID) (bool, error)
	DeleteByVideo(ctx context.Context, videoID valueobject.UUID) error
	DeleteByGenre(ctx context.Context, genreID valueobject.UUID) error
}

// AuditLogRepository is the repository interface for AuditLog
type AuditLogRepository interface {
	Save(ctx context.Context, log *domain.AuditLog) error
	FindByActor(ctx context.Context, actorID valueobject.UUID, limit, offset int) ([]*domain.AuditLog, error)
	FindByResource(ctx context.Context, resourceType, resourceID string, limit, offset int) ([]*domain.AuditLog, error)
	FindRecent(ctx context.Context, limit int) ([]*domain.AuditLog, error)
}

// BatchJobRepository is the repository interface for BatchJob
type BatchJobRepository interface {
	Save(ctx context.Context, job *domain.BatchJob) error
	Update(ctx context.Context, job *domain.BatchJob) error
	FindByID(ctx context.Context, id valueobject.UUID) (*domain.BatchJob, error)
	FindByTypeAndStatus(ctx context.Context, jobType domain.JobType, status domain.JobStatus) ([]*domain.BatchJob, error)
	FindRecent(ctx context.Context, jobType *domain.JobType, limit int) ([]*domain.BatchJob, error)
	GetRunningJobs(ctx context.Context) ([]*domain.BatchJob, error)
}