package input

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// ChannelUseCase defines the interface for channel management use cases
type ChannelUseCase interface {
	ListChannels(ctx context.Context, input ListChannelsInput) (*ListChannelsOutput, error)
	SubscribeChannel(ctx context.Context, youtubeChannelID valueobject.YouTubeChannelID) (*domain.Channel, error)
	UnsubscribeChannel(ctx context.Context, channelID valueobject.UUID) (*domain.Channel, error)
	RenewSubscriptions(ctx context.Context) (*RenewSubscriptionsOutput, error)
}

// ListChannelsInput represents the input for listing channels
type ListChannelsInput struct {
	Subscribed *bool
	Query      *string
	Sort       string
	Limit      int
	Offset     int
}

// ListChannelsOutput represents the output for listing channels
type ListChannelsOutput struct {
	Channels []*domain.Channel
	Total    int
}

// RenewSubscriptionsOutput represents the output for renewing subscriptions
type RenewSubscriptionsOutput struct {
	Renewed int
	Failed  int
}