package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/http/generated"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/presenter"
)

// httpPresenter implements presenter.HTTPPresenter
type httpPresenter struct{}

// NewHTTPPresenter creates a new HTTP presenter
func NewHTTPPresenter() presenter.HTTPPresenter {
	return &httpPresenter{}
}

// Channel operations

func (p *httpPresenter) PresentChannel(channel *domain.Channel) interface{} {
	// Not used in HTTP API - channels are accessed via gRPC
	return &HTTPResponse{
		StatusCode: http.StatusNotImplemented,
		Body: generated.Error{
			Code:    "NOT_IMPLEMENTED",
			Message: "Channel operations are available via gRPC API",
		},
	}
}

func (p *httpPresenter) PresentChannels(channels []*domain.Channel) interface{} {
	// Not used in HTTP API - channels are accessed via gRPC
	return &HTTPResponse{
		StatusCode: http.StatusNotImplemented,
		Body: generated.Error{
			Code:    "NOT_IMPLEMENTED",
			Message: "Channel operations are available via gRPC API",
		},
	}
}

func (p *httpPresenter) PresentUpdateChannelsResult(result interface{}) interface{} {
	r, ok := result.(*UpdateChannelsResult)
	if !ok {
		return p.presentError(http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid result type")
	}

	return &HTTPResponse{
		StatusCode: http.StatusOK,
		Body: generated.UpdateChannelsResponse{
			ChannelsProcessed: int32(r.ChannelsProcessed),
			ChannelsUpdated:   int32(r.ChannelsUpdated),
			Duration:          formatDuration(r.DurationMs),
		},
	}
}

// Video operations

func (p *httpPresenter) PresentVideo(video *domain.Video) interface{} {
	// Not used in HTTP API - videos are accessed via gRPC
	return &HTTPResponse{
		StatusCode: http.StatusNotImplemented,
		Body: generated.Error{
			Code:    "NOT_IMPLEMENTED",
			Message: "Video operations are available via gRPC API",
		},
	}
}

func (p *httpPresenter) PresentVideos(videos []*domain.Video) interface{} {
	// Not used in HTTP API - videos are accessed via gRPC
	return &HTTPResponse{
		StatusCode: http.StatusNotImplemented,
		Body: generated.Error{
			Code:    "NOT_IMPLEMENTED",
			Message: "Video operations are available via gRPC API",
		},
	}
}

func (p *httpPresenter) PresentCollectTrendingResult(result interface{}) interface{} {
	r, ok := result.(*CollectTrendingResult)
	if !ok {
		return p.presentError(http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid result type")
	}

	return &HTTPResponse{
		StatusCode: http.StatusOK,
		Body: generated.CollectTrendingResponse{
			VideosProcessed: int32(r.VideosProcessed),
			VideosAdded:     int32(r.VideosAdded),
			Duration:        formatDuration(r.DurationMs),
		},
	}
}

func (p *httpPresenter) PresentCollectSubscriptionsResult(result interface{}) interface{} {
	r, ok := result.(*CollectSubscriptionsResult)
	if !ok {
		return p.presentError(http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid result type")
	}

	return &HTTPResponse{
		StatusCode: http.StatusOK,
		Body: generated.CollectSubscriptionsResponse{
			ChannelsProcessed: int32(r.ChannelsProcessed),
			VideosProcessed:   int32(r.VideosProcessed),
			VideosAdded:       int32(r.VideosAdded),
			Duration:          formatDuration(r.DurationMs),
		},
	}
}

// Snapshot operations

func (p *httpPresenter) PresentSnapshot(snapshot *domain.VideoSnapshot) interface{} {
	// Not implemented in HTTP API
	return p.presentNotImplemented("Snapshot operations")
}

func (p *httpPresenter) PresentSnapshots(snapshots []*domain.VideoSnapshot) interface{} {
	// Not implemented in HTTP API
	return p.presentNotImplemented("Snapshot operations")
}

// Keyword operations

func (p *httpPresenter) PresentKeyword(keyword *domain.Keyword) interface{} {
	return p.presentNotImplemented("Keyword operations")
}

func (p *httpPresenter) PresentKeywords(keywords []*domain.Keyword) interface{} {
	return p.presentNotImplemented("Keyword operations")
}

func (p *httpPresenter) PresentKeywordCreated(keyword *domain.Keyword) interface{} {
	return p.presentNotImplemented("Keyword operations")
}

func (p *httpPresenter) PresentKeywordUpdated(keyword *domain.Keyword) interface{} {
	return p.presentNotImplemented("Keyword operations")
}

func (p *httpPresenter) PresentKeywordDeleted() interface{} {
	return p.presentNotImplemented("Keyword operations")
}

// Genre operations

func (p *httpPresenter) PresentGenre(genre *domain.Genre) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Genre operations")
}

func (p *httpPresenter) PresentGenres(genres []*domain.Genre) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Genre operations")
}

func (p *httpPresenter) PresentGenreCreated(genre *domain.Genre) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Genre operations")
}

func (p *httpPresenter) PresentGenreUpdated(genre *domain.Genre) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Genre operations")
}

func (p *httpPresenter) PresentGenreEnabled(genre *domain.Genre) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Genre operations")
}

func (p *httpPresenter) PresentGenreDisabled(genre *domain.Genre) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Genre operations")
}

// YouTube Category operations

func (p *httpPresenter) PresentYouTubeCategory(category *domain.YouTubeCategory) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("YouTube Category operations")
}

func (p *httpPresenter) PresentYouTubeCategories(categories []*domain.YouTubeCategory) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("YouTube Category operations")
}

func (p *httpPresenter) PresentYouTubeCategoryCreated(category *domain.YouTubeCategory) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("YouTube Category operations")
}

func (p *httpPresenter) PresentYouTubeCategoryUpdated(category *domain.YouTubeCategory) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("YouTube Category operations")
}

// Video-Genre operations

func (p *httpPresenter) PresentVideoGenre(videoGenre *domain.VideoGenre) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Video-Genre operations")
}

func (p *httpPresenter) PresentVideoGenres(videoGenres []*domain.VideoGenre) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Video-Genre operations")
}

func (p *httpPresenter) PresentVideoGenreCreated(videoGenre *domain.VideoGenre) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Video-Genre operations")
}

func (p *httpPresenter) PresentVideoGenresCreated(videoGenres []*domain.VideoGenre) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Video-Genre operations")
}

func (p *httpPresenter) PresentVideoGenresDeleted() interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Video-Genre operations")
}

// Audit Log operations

func (p *httpPresenter) PresentAuditLog(auditLog *domain.AuditLog) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Audit Log operations")
}

func (p *httpPresenter) PresentAuditLogs(auditLogs []*domain.AuditLog) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Audit Log operations")
}

// Batch Job operations

func (p *httpPresenter) PresentBatchJob(batchJob *domain.BatchJob) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Batch Job operations")
}

func (p *httpPresenter) PresentBatchJobs(batchJobs []*domain.BatchJob) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Batch Job operations")
}

func (p *httpPresenter) PresentBatchJobCreated(batchJob *domain.BatchJob) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Batch Job operations")
}

func (p *httpPresenter) PresentBatchJobUpdated(batchJob *domain.BatchJob) interface{} {
	// Not implemented yet
	return p.presentNotImplemented("Batch Job operations")
}

// System operations

func (p *httpPresenter) PresentScheduleSnapshotsResult(result interface{}) interface{} {
	r, ok := result.(*ScheduleSnapshotsResult)
	if !ok {
		return p.presentError(http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid result type")
	}

	return &HTTPResponse{
		StatusCode: http.StatusOK,
		Body: generated.ScheduleSnapshotsResponse{
			VideosProcessed: int32(r.VideosProcessed),
			TasksScheduled:  int32(r.TasksScheduled),
			Duration:         formatDuration(r.DurationMs),
		},
	}
}

// Common

func (p *httpPresenter) PresentError(err error) *generated.Error {
	switch err {
	case domain.ErrKeywordNotFound:
		return &generated.Error{
			Code:    "NOT_FOUND",
			Message: "Keyword not found",
		}
	case domain.ErrKeywordDuplicate:
		return &generated.Error{
			Code:    "ALREADY_EXISTS",
			Message: "Keyword with the same name already exists",
		}
	case domain.ErrInvalidInput:
		return &generated.Error{
			Code:    "INVALID_ARGUMENT",
			Message: "Invalid input parameters",
		}
	default:
		return &generated.Error{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		}
	}
}

// Helper methods

func (p *httpPresenter) presentError(statusCode int, code, message string) *HTTPResponse {
	return &HTTPResponse{
		StatusCode: statusCode,
		Body: generated.Error{
			Code:    code,
			Message: message,
		},
	}
}

func (p *httpPresenter) presentNotImplemented(operation string) interface{} {
	return &HTTPResponse{
		StatusCode: http.StatusNotImplemented,
		Body: generated.Error{
			Code:    "NOT_IMPLEMENTED",
			Message: operation + " are not available in HTTP API",
		},
	}
}

// formatDuration converts milliseconds to a duration string (e.g., "123ms", "1.5s")
func formatDuration(ms int) string {
	duration := time.Duration(ms) * time.Millisecond
	if duration < time.Second {
		return fmt.Sprintf("%dms", ms)
	}
	return fmt.Sprintf("%.1fs", duration.Seconds())
}