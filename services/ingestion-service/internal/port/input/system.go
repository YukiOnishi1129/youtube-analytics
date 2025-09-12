package input

import (
	"context"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/google/uuid"
)

// SystemInputPort is the interface for system use cases
type SystemInputPort interface {
	ScheduleSnapshots(ctx context.Context) (*ScheduleSnapshotsResult, error)
	CreateSnapshot(ctx context.Context, input *CreateSnapshotInput) (*domain.VideoSnapshot, error)
	GetVideoSnapshots(ctx context.Context, videoID uuid.UUID) ([]*domain.VideoSnapshot, error)
}

// CreateSnapshotInput represents the input for creating a video snapshot
type CreateSnapshotInput struct {
	VideoID        uuid.UUID
	CheckpointHour int
}

// ScheduleSnapshotsResult represents the result of scheduling snapshots
type ScheduleSnapshotsResult struct {
	VideosProcessed int
	TasksScheduled  int
	Duration        time.Duration
}
