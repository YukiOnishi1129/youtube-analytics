package youtube

import (
	"context"
	"fmt"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// client is a real YouTube API client implementation
type client struct {
	service *youtube.Service
	apiKey  string
}

// NewClient creates a new YouTube API client
func NewClient(apiKey string) (gateway.YouTubeClient, error) {
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create youtube service: %w", err)
	}

	return &client{
		service: service,
		apiKey:  apiKey,
	}, nil
}

// GetVideoStats gets video statistics
func (c *client) GetVideoStats(ctx context.Context, ytVideoID valueobject.YouTubeVideoID) (*gateway.VideoStats, error) {
	call := c.service.Videos.List([]string{"statistics"}).Id(string(ytVideoID))
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get video stats: %w", err)
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("video not found: %s", ytVideoID)
	}

	stats := response.Items[0].Statistics
	return &gateway.VideoStats{
		ViewCount:    int64(stats.ViewCount),
		LikeCount:    int64(stats.LikeCount),
		CommentCount: int64(stats.CommentCount),
	}, nil
}

// GetVideoStatistics gets video statistics (string version)
func (c *client) GetVideoStatistics(ctx context.Context, ytVideoID string) (*gateway.VideoStats, error) {
	call := c.service.Videos.List([]string{"statistics"}).Id(ytVideoID)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get video statistics: %w", err)
	}

		if len(response.Items) == 0 {
		return nil, fmt.Errorf("video not found: %s", ytVideoID)
	}

	stats := response.Items[0].Statistics
	return &gateway.VideoStats{
		ViewCount:    int64(stats.ViewCount),
		LikeCount:    int64(stats.LikeCount),
		CommentCount: int64(stats.CommentCount),
	}, nil
}

// GetChannelStats gets channel statistics
func (c *client) GetChannelStats(ctx context.Context, ytChannelID valueobject.YouTubeChannelID) (*gateway.ChannelStats, error) {
	call := c.service.Channels.List([]string{"statistics"}).Id(string(ytChannelID))
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get channel stats: %w", err)
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("channel not found: %s", ytChannelID)
	}

	stats := response.Items[0].Statistics
	return &gateway.ChannelStats{
		SubscriberCount: int64(stats.SubscriberCount),
		VideoCount:      int64(stats.VideoCount),
	}, nil
}

// ListMostPopular lists most popular videos
func (c *client) ListMostPopular(ctx context.Context, categoryID valueobject.CategoryID, pageToken *string) (*gateway.TrendingVideos, error) {
	call := c.service.Videos.List([]string{"snippet"}).
		Chart("mostPopular").
		RegionCode("JP").
		MaxResults(50)

	if categoryID > 0 {
		call = call.VideoCategoryId(fmt.Sprintf("%d", categoryID))
	}

	if pageToken != nil {
		call = call.PageToken(*pageToken)
	}

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list most popular videos: %w", err)
	}

	videos := make([]gateway.VideoMeta, len(response.Items))
	for i, item := range response.Items {
		publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		videos[i] = gateway.VideoMeta{
			ID:           valueobject.YouTubeVideoID(item.Id),
			ChannelID:    valueobject.YouTubeChannelID(item.Snippet.ChannelId),
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishedAt:  publishedAt,
			CategoryID:   categoryID,
			ThumbnailURL: item.Snippet.Thumbnails.High.Url,
		}
	}

	var nextPageToken *string
	if response.NextPageToken != "" {
		nextPageToken = &response.NextPageToken
	}

	return &gateway.TrendingVideos{
		Videos:        videos,
		NextPageToken: nextPageToken,
	}, nil
}

// GetVideo gets video metadata
func (c *client) GetVideo(ctx context.Context, ytVideoID valueobject.YouTubeVideoID) (*gateway.VideoMeta, error) {
	// Implementation using YouTube API
	// This is a placeholder - implement actual API call
	return nil, fmt.Errorf("not implemented")
}

// GetChannel gets channel metadata
func (c *client) GetChannel(ctx context.Context, ytChannelID valueobject.YouTubeChannelID) (*gateway.ChannelMeta, error) {
	call := c.service.Channels.List([]string{"snippet"}).Id(string(ytChannelID))
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("channel not found: %s", ytChannelID)
	}

	channel := response.Items[0]
	return &gateway.ChannelMeta{
		ID:           ytChannelID,
		Title:        channel.Snippet.Title,
		Description:  channel.Snippet.Description,
		ThumbnailURL: channel.Snippet.Thumbnails.High.Url,
	}, nil
}

// GetTrendingVideos gets trending videos
func (c *client) GetTrendingVideos(ctx context.Context) ([]*gateway.VideoMeta, error) {
	trending, err := c.ListMostPopular(ctx, 0, nil)
	if err != nil {
		return nil, err
	}

	result := make([]*gateway.VideoMeta, len(trending.Videos))
	for i, v := range trending.Videos {
		result[i] = &v
	}
	return result, nil
}

// GetChannelVideos gets channel's latest videos
func (c *client) GetChannelVideos(ctx context.Context, channelID valueobject.YouTubeChannelID) ([]*gateway.VideoMeta, error) {
	call := c.service.Search.List([]string{"snippet"}).
		ChannelId(string(channelID)).
		Order("date").
		MaxResults(50)

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get channel videos: %w", err)
	}

	videos := make([]*gateway.VideoMeta, 0, len(response.Items))
	for _, item := range response.Items {
		if item.Id.Kind != "youtube#video" {
			continue
		}

		publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		videos = append(videos, &gateway.VideoMeta{
			ID:           valueobject.YouTubeVideoID(item.Id.VideoId),
			ChannelID:    channelID,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishedAt:  publishedAt,
			ThumbnailURL: item.Snippet.Thumbnails.High.Url,
		})
	}

	return videos, nil
}
