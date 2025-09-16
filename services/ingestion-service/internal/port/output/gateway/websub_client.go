package gateway

import (
	"context"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// WebSubHub is the gateway interface for WebSub Hub operations
type WebSubHub interface {
	Subscribe(ctx context.Context, channelID valueobject.YouTubeChannelID, callbackURL string, leaseSeconds int) error
	Unsubscribe(ctx context.Context, channelID valueobject.YouTubeChannelID, callbackURL string) error
}

// WebSubNotification represents a notification from WebSub hub
type WebSubNotification struct {
	VideoID      valueobject.YouTubeVideoID
	ChannelID    valueobject.YouTubeChannelID
	PublishedAt  time.Time
	ReceivedAt   time.Time
}