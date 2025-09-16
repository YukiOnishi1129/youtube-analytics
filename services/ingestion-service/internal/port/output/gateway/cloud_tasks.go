package gateway

import (
	"context"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// CloudTasksGateway is the interface for Google Cloud Tasks
type CloudTasksGateway interface {
	CreateSnapshotTask(ctx context.Context, videoID valueobject.UUID, checkpointHour valueobject.CheckpointHour, scheduleTime time.Time) error
	DeleteSnapshotTask(ctx context.Context, videoID valueobject.UUID, checkpointHour valueobject.CheckpointHour) error
}