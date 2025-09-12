package input

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
)

// ChannelInputPort is the interface for channel use cases
type ChannelInputPort interface {
	UpdateChannels(ctx context.Context) (*UpdateChannelsResult, error)
	GetChannel(ctx context.Context, channelID uuid.UUID) (*domain.Channel, error)
	ListChannels(ctx context.Context, onlySubscribed bool) ([]*domain.Channel, error)
}

// UpdateChannelsResult represents the result of updating channels
type UpdateChannelsResult struct {
	ChannelsProcessed int
	ChannelsUpdated   int
	Duration          time.Duration
}