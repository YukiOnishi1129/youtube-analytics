package domain

import (
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// SnapshotTask represents a scheduled task to capture video snapshot
type SnapshotTask struct {
	VideoID        valueobject.UUID
	CheckpointHour valueobject.CheckpointHour
	ScheduledAt    time.Time
}