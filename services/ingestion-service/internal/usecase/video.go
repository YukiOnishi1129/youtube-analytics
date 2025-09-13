package usecase

import (
	"context"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

type videoUseCase struct {
	videoRepo      gateway.VideoRepository
	channelRepo    gateway.ChannelRepository
	youtubeAPI     gateway.YouTubeClient
	eventPublisher gateway.EventPublisher
}

func NewVideoUseCase(
	videoRepo gateway.VideoRepository,
	channelRepo gateway.ChannelRepository,
	youtubeAPI gateway.YouTubeClient,
	eventPublisher gateway.EventPublisher,
) input.VideoInputPort {
	return &videoUseCase{
		videoRepo:      videoRepo,
		channelRepo:    channelRepo,
		youtubeAPI:     youtubeAPI,
		eventPublisher: eventPublisher,
	}
}

func (u *videoUseCase) CollectTrending(ctx context.Context) (*input.CollectTrendingResult, error) {
	start := time.Now()
	// Fetch trending videos from YouTube API
	trendingVideos, err := u.youtubeAPI.GetTrendingVideos(ctx)
	if err != nil {
		return nil, err
	}

	videosAdded := 0
	for _, videoMeta := range trendingVideos {
		// Check if video already exists
		exists, err := u.videoRepo.ExistsByYouTubeVideoID(ctx, videoMeta.ID)
		if err != nil {
			continue
		}
		if exists {
			continue
		}

		// Create new video
		video := &domain.Video{
			ID:               valueobject.GenerateUUID(),
			YouTubeVideoID:   valueobject.YouTubeVideoID(videoMeta.ID),
			YouTubeChannelID: valueobject.YouTubeChannelID(videoMeta.ChannelID),
			Title:            videoMeta.Title,
			ThumbnailURL:     videoMeta.ThumbnailURL,
			PublishedAt:      videoMeta.PublishedAt,
			CreatedAt:        time.Now(),
		}

		if err := u.videoRepo.Save(ctx, video); err != nil {
			continue
		}

		// Publish event
		if err := u.eventPublisher.PublishVideoDiscovered(ctx, video); err != nil {
			// Log error but continue
			continue
		}

		videosAdded++
	}

	return &input.CollectTrendingResult{
		VideosProcessed: len(trendingVideos),
		VideosAdded:     videosAdded,
		Duration:        time.Since(start),
	}, nil
}

func (u *videoUseCase) CollectSubscriptions(ctx context.Context) (*input.CollectSubscriptionsResult, error) {
	start := time.Now()
	// Get all subscribed channels
	channels, err := u.channelRepo.ListSubscribed(ctx)
	if err != nil {
		return nil, err
	}

	videosAdded := 0
	totalVideos := 0

	for _, channel := range channels {
		// Fetch latest videos from channel
		videos, err := u.youtubeAPI.GetChannelVideos(ctx, channel.YouTubeChannelID)
		if err != nil {
			continue
		}

		totalVideos += len(videos)

		for _, videoMeta := range videos {
			// Check if video already exists
			exists, err := u.videoRepo.ExistsByYouTubeVideoID(ctx, videoMeta.ID)
			if err != nil {
				continue
			}
			if exists {
				continue
			}

			// Create new video
			video := &domain.Video{
				ID:               valueobject.GenerateUUID(),
				YouTubeVideoID:   valueobject.YouTubeVideoID(videoMeta.ID),
				YouTubeChannelID: channel.YouTubeChannelID,
				Title:            videoMeta.Title,
				ThumbnailURL:     videoMeta.ThumbnailURL,
				PublishedAt:      videoMeta.PublishedAt,
				CreatedAt:        time.Now(),
			}

			if err := u.videoRepo.Save(ctx, video); err != nil {
				continue
			}

			// Publish event
			if err := u.eventPublisher.PublishVideoDiscovered(ctx, video); err != nil {
				// Log error but continue
				continue
			}

			videosAdded++
		}
	}

	return &input.CollectSubscriptionsResult{
		ChannelsProcessed: len(channels),
		VideosProcessed:   totalVideos,
		VideosAdded:       videosAdded,
		Duration:          time.Since(start),
	}, nil
}
