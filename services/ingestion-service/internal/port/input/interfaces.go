package input

import "context"

// ChannelInputPort is the interface for channel use cases
type ChannelInputPort interface {
	UpdateChannels(ctx context.Context) (*UpdateChannelsResult, error)
}

// VideoInputPort is the interface for video use cases
type VideoInputPort interface {
	CollectTrending(ctx context.Context) (*CollectTrendingResult, error)
	CollectSubscriptions(ctx context.Context) (*CollectSubscriptionsResult, error)
}

// SystemInputPort is the interface for system use cases
type SystemInputPort interface {
	ScheduleSnapshots(ctx context.Context) (*ScheduleSnapshotsResult, error)
}

// Result types for ChannelInputPort
type UpdateChannelsResult struct {
	ChannelsProcessed int
	ChannelsUpdated   int
}

// Result types for VideoInputPort
type CollectTrendingResult struct {
	VideosProcessed int
	VideosAdded     int
}

type CollectSubscriptionsResult struct {
	ChannelsProcessed int
	VideosProcessed   int
	VideosAdded       int
}

// Result types for SystemInputPort
type ScheduleSnapshotsResult struct {
	VideosProcessed int
	TasksScheduled  int
}