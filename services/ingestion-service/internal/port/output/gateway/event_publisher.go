package gateway

import (
	"context"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// EventPublisher is the gateway interface for publishing domain events
type EventPublisher interface {
	PublishSnapshotAdded(ctx context.Context, event SnapshotAddedEvent) error
}

// SnapshotAddedEvent represents a domain event for snapshot creation
type SnapshotAddedEvent struct {
	VideoID        valueobject.UUID
	CheckpointHour valueobject.CheckpointHour
	MeasuredAt     time.Time
}