package usecase

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
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
	}, nil
}