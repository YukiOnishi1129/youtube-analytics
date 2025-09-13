package input

import (
	"context"
	"time"
)

// VideoInputPort is the interface for video use cases
type VideoInputPort interface {
	CollectTrending(ctx context.Context) (*CollectTrendingResult, error)
	CollectSubscriptions(ctx context.Context) (*CollectSubscriptionsResult, error)
}

// CollectTrendingResult represents the result of collecting trending videos
type CollectTrendingResult struct {
	VideosProcessed int
	VideosAdded     int
	Duration        time.Duration
}

// CollectSubscriptionsResult represents the result of collecting subscription videos
type CollectSubscriptionsResult struct {
	ChannelsProcessed int
	VideosProcessed   int
	VideosAdded       int
	Duration          time.Duration
}