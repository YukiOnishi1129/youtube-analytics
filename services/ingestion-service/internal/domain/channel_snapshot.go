package domain

import (
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// ChannelSnapshot represents a snapshot of channel statistics
type ChannelSnapshot struct {
	ID                valueobject.UUID
	ChannelID         valueobject.UUID
	MeasuredAt        time.Time
	SubscriptionCount int
	CreatedAt         time.Time
}

// NewChannelSnapshot creates a new channel snapshot
func NewChannelSnapshot(
	id valueobject.UUID,
	channelID valueobject.UUID,
	measuredAt time.Time,
	subscriptionCount int,
) *ChannelSnapshot {
	return &ChannelSnapshot{
		ID:                id,
		ChannelID:         channelID,
		MeasuredAt:        measuredAt,
		SubscriptionCount: subscriptionCount,
		CreatedAt:         time.Now(),
	}
}