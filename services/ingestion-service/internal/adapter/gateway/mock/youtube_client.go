package mock

import (
	"context"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

// youtubeClient is a mock implementation of YouTubeClient
type youtubeClient struct{}

// NewYouTubeClient creates a new mock YouTube client
func NewYouTubeClient() gateway.YouTubeClient {
	return &youtubeClient{}
}

func (c *youtubeClient) GetVideoStats(ctx context.Context, ytVideoID valueobject.YouTubeVideoID) (*gateway.VideoStats, error) {
	return &gateway.VideoStats{
		ViewCount:    1000,
		LikeCount:    100,
		CommentCount: 50,
	}, nil
}

func (c *youtubeClient) GetChannelStats(ctx context.Context, ytChannelID valueobject.YouTubeChannelID) (*gateway.ChannelStats, error) {
	return &gateway.ChannelStats{
		SubscriberCount: 10000,
		VideoCount:      100,
	}, nil
}

func (c *youtubeClient) ListMostPopular(ctx context.Context, categoryID valueobject.CategoryID, pageToken *string) (*gateway.TrendingVideos, error) {
	return &gateway.TrendingVideos{
		Videos: []gateway.VideoMeta{
			{
				ID:           "video1",
				ChannelID:    "channel1",
				Title:        "Trending Video 1",
				Description:  "Description 1",
				PublishedAt:  time.Now().Add(-24 * time.Hour),
				CategoryID:   categoryID,
				ThumbnailURL: "https://example.com/thumb1.jpg",
			},
			{
				ID:           "video2",
				ChannelID:    "channel2",
				Title:        "Trending Video 2",
				Description:  "Description 2",
				PublishedAt:  time.Now().Add(-48 * time.Hour),
				CategoryID:   categoryID,
				ThumbnailURL: "https://example.com/thumb2.jpg",
			},
		},
		NextPageToken: nil,
	}, nil
}

func (c *youtubeClient) GetVideo(ctx context.Context, ytVideoID valueobject.YouTubeVideoID) (*gateway.VideoMeta, error) {
	return &gateway.VideoMeta{
		ID:           ytVideoID,
		ChannelID:    "channel1",
		Title:        "Video Title",
		Description:  "Video Description",
		PublishedAt:  time.Now().Add(-24 * time.Hour),
		CategoryID:   10,
		ThumbnailURL: "https://example.com/thumb.jpg",
	}, nil
}

func (c *youtubeClient) GetChannel(ctx context.Context, ytChannelID valueobject.YouTubeChannelID) (*gateway.ChannelMeta, error) {
	return &gateway.ChannelMeta{
		ID:           ytChannelID,
		Title:        "Channel Title",
		Description:  "Channel Description",
		ThumbnailURL: "https://example.com/channel.jpg",
	}, nil
}

func (c *youtubeClient) GetTrendingVideos(ctx context.Context) ([]*gateway.VideoMeta, error) {
	return []*gateway.VideoMeta{
		{
			ID:           "trending1",
			ChannelID:    "channel1",
			Title:        "Trending Video 1",
			Description:  "Description 1",
			PublishedAt:  time.Now().Add(-24 * time.Hour),
			CategoryID:   10,
			ThumbnailURL: "https://example.com/trending1.jpg",
		},
		{
			ID:           "trending2",
			ChannelID:    "channel2",
			Title:        "Trending Video 2",
			Description:  "Description 2",
			PublishedAt:  time.Now().Add(-48 * time.Hour),
			CategoryID:   10,
			ThumbnailURL: "https://example.com/trending2.jpg",
		},
	}, nil
}

func (c *youtubeClient) GetChannelVideos(ctx context.Context, channelID valueobject.YouTubeChannelID) ([]*gateway.VideoMeta, error) {
	return []*gateway.VideoMeta{
		{
			ID:           "channelvideo1",
			ChannelID:    channelID,
			Title:        "Channel Video 1",
			Description:  "Description 1",
			PublishedAt:  time.Now().Add(-24 * time.Hour),
			CategoryID:   10,
			ThumbnailURL: "https://example.com/channelvideo1.jpg",
		},
		{
			ID:           "channelvideo2",
			ChannelID:    channelID,
			Title:        "Channel Video 2",
			Description:  "Description 2",
			PublishedAt:  time.Now().Add(-48 * time.Hour),
			CategoryID:   10,
			ThumbnailURL: "https://example.com/channelvideo2.jpg",
		},
	}, nil
}