package input

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
)

// SystemInputPort is the interface for system use cases
type SystemInputPort interface {
	ScheduleSnapshots(ctx context.Context) (*ScheduleSnapshotsResult, error)
	CreateSnapshot(ctx context.Context, videoID uuid.UUID, checkpointHour int) (*domain.VideoSnapshot, error)
	GetVideoSnapshots(ctx context.Context, videoID uuid.UUID) ([]*domain.VideoSnapshot, error)
}

// ScheduleSnapshotsResult represents the result of scheduling snapshots
type ScheduleSnapshotsResult struct {
	VideosProcessed int
	TasksScheduled  int
	Duration        time.Duration
}