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
	VideosCollected int // Total videos fetched from YouTube API
	VideosCreated   int // New videos added to database
	VideosUpdated   int // Existing videos updated
	Duration        time.Duration
}

// CollectAllTrendingResult represents the result of collecting trending videos for all genres
type CollectAllTrendingResult struct {
	GenresProcessed int
	TotalCollected  int // Total videos collected from all genres
	TotalCreated    int // Total new videos created
	TotalUpdated    int // Total videos updated
	GenreResults    []*CollectTrendingResult
	Duration        time.Duration
}

// CollectSubscriptionsResult represents the result of collecting subscription videos
type CollectSubscriptionsResult struct {
	ChannelsProcessed int
	VideosCollected   int // Total videos fetched from subscribed channels
	VideosCreated     int // New videos added to database
	Duration          time.Duration
}