package gateway

import (
	"context"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// YouTubeClient is the gateway interface for YouTube Data API
type YouTubeClient interface {
	GetVideoStats(ctx context.Context, ytVideoID valueobject.YouTubeVideoID) (*VideoStats, error)
	GetChannelStats(ctx context.Context, ytChannelID valueobject.YouTubeChannelID) (*ChannelStats, error)
	ListMostPopular(ctx context.Context, categoryID valueobject.CategoryID, pageToken *string) (*TrendingVideos, error)
	GetVideo(ctx context.Context, ytVideoID valueobject.YouTubeVideoID) (*VideoMeta, error)
	GetChannel(ctx context.Context, ytChannelID valueobject.YouTubeChannelID) (*ChannelMeta, error)
}

// VideoStats represents video statistics from YouTube API
type VideoStats struct {
	ViewCount    int64
	LikeCount    int64
	CommentCount int64
}

// ChannelStats represents channel statistics from YouTube API
type ChannelStats struct {
	SubscriberCount int64
	VideoCount      int64
}

// VideoMeta represents video metadata from YouTube API
type VideoMeta struct {
	ID          valueobject.YouTubeVideoID
	ChannelID   valueobject.YouTubeChannelID
	Title       string
	Description string
	PublishedAt time.Time
	CategoryID  valueobject.CategoryID
	Thumbnails  Thumbnails
}

// ChannelMeta represents channel metadata from YouTube API
type ChannelMeta struct {
	ID          valueobject.YouTubeChannelID
	Title       string
	Description string
	Thumbnails  Thumbnails
}

// Thumbnails represents YouTube thumbnails
type Thumbnails struct {
	Default  string
	Medium   string
	High     string
	Standard string
	MaxRes   string
}

// TrendingVideos represents a page of trending videos
type TrendingVideos struct {
	Videos        []VideoMeta
	NextPageToken *string
}

