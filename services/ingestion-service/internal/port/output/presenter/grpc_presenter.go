package presenter

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
)

// GRPCPresenter is the output port for gRPC presentation
type GRPCPresenter interface {
	// Channel operations
	PresentChannel(channel *domain.Channel) interface{}
	PresentChannels(channels []*domain.Channel) interface{}
	
	// Video operations
	PresentVideo(video *domain.Video) interface{}
	PresentVideos(videos []*domain.Video) interface{}
	
	// Snapshot operations
	PresentSnapshot(snapshot *domain.VideoSnapshot) interface{}
	PresentSnapshots(snapshots []*domain.VideoSnapshot) interface{}
	
	// Keyword operations
	PresentKeyword(keyword *domain.Keyword) interface{}
	PresentKeywords(keywords []*domain.Keyword) interface{}
	
	// System operations
	PresentScheduleSnapshotsResult(result interface{}) interface{}
	PresentUpdateChannelsResult(result interface{}) interface{}
	PresentCollectTrendingResult(result interface{}) interface{}
	PresentCollectSubscriptionsResult(result interface{}) interface{}
	
	// Common
	PresentDeleted() interface{}
	PresentError(err error) error
}