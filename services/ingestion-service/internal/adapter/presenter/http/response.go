package http

// HTTPResponse represents a standard HTTP response
type HTTPResponse struct {
	StatusCode int
	Body       interface{}
	Headers    map[string]string
}

// Result types for system operations

// ScheduleSnapshotsResult represents the result of scheduling snapshots
type ScheduleSnapshotsResult struct {
	VideosProcessed int
	TasksScheduled  int
	DurationMs      int
}

// UpdateChannelsResult represents the result of updating channels
type UpdateChannelsResult struct {
	ChannelsProcessed int
	ChannelsUpdated   int
	DurationMs        int
}

// CollectTrendingResult represents the result of collecting trending videos
type CollectTrendingResult struct {
	VideosProcessed int
	VideosAdded     int
	DurationMs      int
}

// CollectSubscriptionsResult represents the result of collecting subscription videos
type CollectSubscriptionsResult struct {
	ChannelsProcessed int
	VideosProcessed   int
	VideosAdded       int
	DurationMs        int
}