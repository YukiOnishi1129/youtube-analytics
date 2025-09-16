package usecase

import (
	"context"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"github.com/google/uuid"
)

type channelUseCase struct {
	channelRepo gateway.ChannelRepository
	youtubeAPI  gateway.YouTubeClient
}

func NewChannelUseCase(
	channelRepo gateway.ChannelRepository,
	youtubeAPI gateway.YouTubeClient,
) input.ChannelInputPort {
	return &channelUseCase{
		channelRepo: channelRepo,
		youtubeAPI:  youtubeAPI,
	}
}

func (u *channelUseCase) UpdateChannels(ctx context.Context) (*input.UpdateChannelsResult, error) {
	start := time.Now()

	// Get all active channels
	channels, err := u.channelRepo.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	updated := 0
	for _, channel := range channels {
		// Fetch latest metadata from YouTube API
		metadata, err := u.youtubeAPI.GetChannel(ctx, channel.YouTubeChannelID)
		if err != nil {
			// Continue with next channel on error
			continue
		}

		// Update if changed
		if channel.Title != metadata.Title || channel.ThumbnailURL != metadata.ThumbnailURL {
			channel.Title = metadata.Title
			channel.ThumbnailURL = metadata.ThumbnailURL

			if err := u.channelRepo.Update(ctx, channel); err != nil {
				// Continue with next channel on error
				continue
			}
			updated++
		}
	}

	return &input.UpdateChannelsResult{
		ChannelsProcessed: len(channels),
		ChannelsUpdated:   updated,
		Duration:          time.Since(start),
	}, nil
}

func (u *channelUseCase) GetChannel(ctx context.Context, channelID uuid.UUID) (*domain.Channel, error) {
	channel, err := u.channelRepo.GetByID(ctx, valueobject.UUID(channelID.String()))
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (u *channelUseCase) ListChannels(ctx context.Context, onlySubscribed bool) ([]*domain.Channel, error) {
	if onlySubscribed {
		return u.channelRepo.ListSubscribed(ctx)
	}
	return u.channelRepo.ListActive(ctx)
}
