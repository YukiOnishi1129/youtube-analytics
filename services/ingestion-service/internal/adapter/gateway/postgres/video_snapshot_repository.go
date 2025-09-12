package postgres

import (
	"context"
	"database/sql"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

// videoSnapshotRepository implements gateway.VideoSnapshotRepository interface
type videoSnapshotRepository struct {
	*Repository
}

// NewVideoSnapshotRepository creates a new video snapshot repository
func NewVideoSnapshotRepository(repo *Repository) gateway.VideoSnapshotRepository {
	return &videoSnapshotRepository{Repository: repo}
}


// Exists checks if a snapshot exists for a video at a specific checkpoint
func (r *videoSnapshotRepository) Exists(ctx context.Context, videoID valueobject.UUID, cp valueobject.CheckpointHour) (bool, error) {
	// TODO: Implement once SQL query is generated
	// For now, return false
	return false, nil
}

// FindByVideoAndCP finds a snapshot by video ID and checkpoint
func (r *videoSnapshotRepository) FindByVideoAndCP(ctx context.Context, videoID valueobject.UUID, cp valueobject.CheckpointHour) (*domain.VideoSnapshot, error) {
	// TODO: Implement once SQL query is generated
	// For now, return not found
	return nil, sql.ErrNoRows
}

// ListByVideo lists all snapshots for a video
func (r *videoSnapshotRepository) ListByVideo(ctx context.Context, videoID valueobject.UUID) ([]*domain.VideoSnapshot, error) {
	return r.ListByVideoID(ctx, videoID)
}

// ListByVideoID lists all snapshots for a video ID
func (r *videoSnapshotRepository) ListByVideoID(ctx context.Context, videoID valueobject.UUID) ([]*domain.VideoSnapshot, error) {
	// TODO: Implement once SQL query is generated
	// For now, return empty list
	return []*domain.VideoSnapshot{}, nil
}
