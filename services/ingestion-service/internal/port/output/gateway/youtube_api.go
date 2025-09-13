package gateway

import (
	"context"
)

// YouTubeAPI is the legacy interface for YouTube Data API (to be replaced by YouTubeClient)
type YouTubeAPI interface {
	GetChannel(ctx context.Context, channelID string) (*ChannelMeta, error)
	GetVideo(ctx context.Context, videoID string) (*VideoMeta, error)
}
