package input

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

// VideoUseCase defines the interface for video monitoring and snapshot use cases
type VideoUseCase interface {
	RegisterFromTrending(ctx context.Context, input RegisterFromTrendingInput) (*RegisterFromTrendingOutput, error)
	ApplyWebSubNotification(ctx context.Context, youtubeVideoID valueobject.YouTubeVideoID) error
	AddSnapshot(ctx context.Context, input AddSnapshotInput) error
}

// RegisterFromTrendingInput represents the input for registering videos from trending
type RegisterFromTrendingInput struct {
	Videos []gateway.VideoMeta
}

// RegisterFromTrendingOutput represents the output for registering videos from trending
type RegisterFromTrendingOutput struct {
	Registered int
	Excluded   int
	Neutral    int
}

// AddSnapshotInput represents the input for adding a snapshot
type AddSnapshotInput struct {
	VideoID        valueobject.UUID
	CheckpointHour valueobject.CheckpointHour
}