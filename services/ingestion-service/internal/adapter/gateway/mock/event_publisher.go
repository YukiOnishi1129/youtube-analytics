package mock

import (
	"context"
	"log"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

// eventPublisher is a mock implementation of EventPublisher
type eventPublisher struct{}

// NewEventPublisher creates a new mock event publisher
func NewEventPublisher() gateway.EventPublisher {
	return &eventPublisher{}
}

// PublishSnapshotAdded publishes a snapshot added event
func (p *eventPublisher) PublishSnapshotAdded(ctx context.Context, event gateway.SnapshotAddedEvent) error {
	log.Printf("Mock: Publishing snapshot added event for video %s at checkpoint %d", event.VideoID, event.CheckpointHour)
	return nil
}

// PublishVideoDiscovered publishes a video discovered event
func (p *eventPublisher) PublishVideoDiscovered(ctx context.Context, video *domain.Video) error {
	log.Printf("Mock: Publishing video discovered event for video %s", video.ID)
	return nil
}