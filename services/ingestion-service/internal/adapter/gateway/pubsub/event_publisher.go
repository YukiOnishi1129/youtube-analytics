package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

// eventPublisher is a Google Cloud Pub/Sub implementation of EventPublisher
type eventPublisher struct {
	client               *pubsub.Client
	snapshotAddedTopic   *pubsub.Topic
	videoDiscoveredTopic *pubsub.Topic
}

// NewEventPublisher creates a new Pub/Sub event publisher
func NewEventPublisher(projectID string) (gateway.EventPublisher, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}

	return &eventPublisher{
		client:               client,
		snapshotAddedTopic:   client.Topic("snapshot-added"),
		videoDiscoveredTopic: client.Topic("video-discovered"),
	}, nil
}

// PublishSnapshotAdded publishes a snapshot added event
func (p *eventPublisher) PublishSnapshotAdded(ctx context.Context, event gateway.SnapshotAddedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	result := p.snapshotAddedTopic.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	// Wait for the publish to complete
	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// PublishVideoDiscovered publishes a video discovered event
func (p *eventPublisher) PublishVideoDiscovered(ctx context.Context, video *domain.Video) error {
	event := map[string]interface{}{
		"videoId":          string(video.ID),
		"youtubeVideoId":   string(video.YouTubeVideoID),
		"youtubeChannelId": string(video.YouTubeChannelID),
		"title":            video.Title,
		"publishedAt":      video.PublishedAt,
		"createdAt":        video.CreatedAt,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	result := p.videoDiscoveredTopic.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	// Wait for the publish to complete
	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// Close closes the pubsub client
func (p *eventPublisher) Close() error {
	return p.client.Close()
}
