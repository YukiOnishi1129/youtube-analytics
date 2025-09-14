package input

import (
	"context"
	"time"
)

// VideoInputPort is the interface for video use cases
type VideoInputPort interface {
	CollectTrending(ctx context.Context, genreID *string) (*CollectTrendingResult, error)
	CollectTrendingByGenre(ctx context.Context, genreID string) (*CollectTrendingResult, error)
	CollectAllTrending(ctx context.Context) (*CollectAllTrendingResult, error)
	CollectSubscriptions(ctx context.Context) (*CollectSubscriptionsResult, error)
}

// CollectTrendingResult represents the result of collecting trending videos
type CollectTrendingResult struct {
	GenreCode       string
	VideosProcessed int
	VideosAdded     int
	Duration        time.Duration
}

// CollectAllTrendingResult represents the result of collecting trending videos for all genres
type CollectAllTrendingResult struct {
	GenresProcessed int
	TotalVideos     int
	TotalAdded      int
	GenreResults    []*CollectTrendingResult
	Duration        time.Duration
}

// CollectSubscriptionsResult represents the result of collecting subscription videos
type CollectSubscriptionsResult struct {
	ChannelsProcessed int
	VideosProcessed   int
	VideosAdded       int
	Duration          time.Duration
}