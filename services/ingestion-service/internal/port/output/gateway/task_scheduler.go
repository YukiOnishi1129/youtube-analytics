package gateway

import (
	"context"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// TaskScheduler is the gateway interface for Cloud Tasks
type TaskScheduler interface {
	Schedule(ctx context.Context, videoID valueobject.UUID, cp valueobject.CheckpointHour, eta time.Time) error
	Cancel(ctx context.Context, videoID valueobject.UUID, cp valueobject.CheckpointHour) error
	ScheduleSnapshot(ctx context.Context, task *domain.SnapshotTask) error
}