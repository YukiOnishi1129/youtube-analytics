package presenter

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/http/generated"
)

// HTTPPresenter is the output port for HTTP presentation
type HTTPPresenter interface {
	// Channel operations
	PresentChannel(channel *domain.Channel) interface{}
	PresentChannels(channels []*domain.Channel) interface{}
	PresentUpdateChannelsResult(result interface{}) interface{}
	
	// Video operations
	PresentVideo(video *domain.Video) interface{}
	PresentVideos(videos []*domain.Video) interface{}
	PresentCollectTrendingResult(result interface{}) interface{}
	PresentCollectSubscriptionsResult(result interface{}) interface{}
	
	// Snapshot operations
	PresentSnapshot(snapshot *domain.VideoSnapshot) interface{}
	PresentSnapshots(snapshots []*domain.VideoSnapshot) interface{}
	
	// Keyword operations
	PresentKeyword(keyword *domain.Keyword) interface{}
	PresentKeywords(keywords []*domain.Keyword) interface{}
	PresentKeywordCreated(keyword *domain.Keyword) interface{}
	PresentKeywordUpdated(keyword *domain.Keyword) interface{}
	PresentKeywordDeleted() interface{}
	
	// Genre operations
	PresentGenre(genre *domain.Genre) interface{}
	PresentGenres(genres []*domain.Genre) interface{}
	PresentGenreCreated(genre *domain.Genre) interface{}
	PresentGenreUpdated(genre *domain.Genre) interface{}
	PresentGenreEnabled(genre *domain.Genre) interface{}
	PresentGenreDisabled(genre *domain.Genre) interface{}
	
	// YouTube Category operations
	PresentYouTubeCategory(category *domain.YouTubeCategory) interface{}
	PresentYouTubeCategories(categories []*domain.YouTubeCategory) interface{}
	PresentYouTubeCategoryCreated(category *domain.YouTubeCategory) interface{}
	PresentYouTubeCategoryUpdated(category *domain.YouTubeCategory) interface{}
	
	// Video-Genre operations
	PresentVideoGenre(videoGenre *domain.VideoGenre) interface{}
	PresentVideoGenres(videoGenres []*domain.VideoGenre) interface{}
	PresentVideoGenreCreated(videoGenre *domain.VideoGenre) interface{}
	PresentVideoGenresCreated(videoGenres []*domain.VideoGenre) interface{}
	PresentVideoGenresDeleted() interface{}
	
	// Audit Log operations
	PresentAuditLog(auditLog *domain.AuditLog) interface{}
	PresentAuditLogs(auditLogs []*domain.AuditLog) interface{}
	
	// Batch Job operations
	PresentBatchJob(batchJob *domain.BatchJob) interface{}
	PresentBatchJobs(batchJobs []*domain.BatchJob) interface{}
	PresentBatchJobCreated(batchJob *domain.BatchJob) interface{}
	PresentBatchJobUpdated(batchJob *domain.BatchJob) interface{}
	
	// System operations
	PresentScheduleSnapshotsResult(result interface{}) interface{}
	
	// Common
	PresentError(err error) *generated.Error
}