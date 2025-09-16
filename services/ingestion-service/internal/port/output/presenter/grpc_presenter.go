package presenter

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	pb "github.com/YukiOnishi1129/youtube-analytics/services/pkg/pb/ingestion/v1"
)

// GRPCPresenter is the output port for gRPC presentation
type GRPCPresenter interface {
	// Channel operations
	PresentChannel(channel *domain.Channel) *pb.Channel
	PresentChannels(channels []*domain.Channel) []*pb.Channel

	// Video operations
	PresentVideo(video *domain.Video) *pb.Video
	PresentVideos(videos []*domain.Video) []*pb.Video

	// Snapshot operations
	PresentSnapshot(snapshot *domain.VideoSnapshot) *pb.VideoSnapshot
	PresentSnapshots(snapshots []*domain.VideoSnapshot) []*pb.VideoSnapshot

	// Keyword operations
	PresentKeyword(keyword *domain.Keyword) *pb.Keyword
	PresentKeywords(keywords []*domain.Keyword) []*pb.Keyword

	// Genre operations
	PresentGenre(genre *domain.Genre) *pb.Genre
	PresentGenres(genres []*domain.Genre) []*pb.Genre

	// YouTube Category operations
	PresentYouTubeCategory(category *domain.YouTubeCategory) *pb.YouTubeCategory
	PresentYouTubeCategories(categories []*domain.YouTubeCategory) []*pb.YouTubeCategory

	// Video-Genre operations
	PresentVideoGenre(videoGenre *domain.VideoGenre) *pb.VideoGenre
	PresentVideoGenres(videoGenres []*domain.VideoGenre) []*pb.VideoGenre

	// Audit Log operations
	PresentAuditLog(auditLog *domain.AuditLog) *pb.AuditLog
	PresentAuditLogs(auditLogs []*domain.AuditLog) []*pb.AuditLog

	// Batch Job operations
	PresentBatchJob(batchJob *domain.BatchJob) *pb.BatchJob
	PresentBatchJobs(batchJobs []*domain.BatchJob) []*pb.BatchJob

	// System operations
	PresentScheduleSnapshotsResult(count int32) *pb.ScheduleSnapshotsResponse
	PresentUpdateChannelsResult(count int32) *pb.UpdateChannelsResponse
	PresentCollectTrendingResult(count int32) *pb.CollectTrendingResponse
	PresentCollectSubscriptionsResult(count int32) *pb.CollectSubscriptionsResponse

	// Common
	PresentError(err error) error
}
