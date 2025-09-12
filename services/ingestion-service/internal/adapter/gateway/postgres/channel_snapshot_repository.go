package postgres

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

// channelSnapshotRepository implements gateway.ChannelSnapshotRepository interface
type channelSnapshotRepository struct {
	*Repository
}

// NewChannelSnapshotRepository creates a new channel snapshot repository
func NewChannelSnapshotRepository(repo *Repository) gateway.ChannelSnapshotRepository {
	return &channelSnapshotRepository{Repository: repo}
}


// Latest gets the latest snapshot for a channel
func (r *channelSnapshotRepository) Latest(ctx context.Context, channelID valueobject.UUID) (*domain.ChannelSnapshot, error) {
	// TODO: Implement get latest channel snapshot
	return nil, nil
}

// ListByChannel lists snapshots for a channel
func (r *channelSnapshotRepository) ListByChannel(ctx context.Context, channelID valueobject.UUID, limit int) ([]*domain.ChannelSnapshot, error) {
	// TODO: Implement list channel snapshots
	return nil, nil
}
