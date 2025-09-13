package http

import (
	"net/http"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/http/generated"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/presenter"
)

// videoPresenter implements presenter.VideoPresenter for HTTP
type videoPresenter struct{}

// NewVideoPresenter creates a new HTTP video presenter
func NewVideoPresenter() presenter.VideoPresenter {
	return &videoPresenter{}
}

// PresentVideo presents a single video
// Note: This is for HTTP REST API, not used in current implementation
// as video operations are exposed via gRPC only
func (p *videoPresenter) PresentVideo(video *domain.Video) interface{} {
	// Not used in HTTP API - videos are accessed via gRPC
	return &HTTPResponse{
		StatusCode: http.StatusNotImplemented,
		Body: generated.Error{
			Code:    "NOT_IMPLEMENTED",
			Message: "Video operations are available via gRPC API",
		},
	}
}

// PresentVideos presents multiple videos  
// Note: This is for HTTP REST API, not used in current implementation
// as video operations are exposed via gRPC only
func (p *videoPresenter) PresentVideos(videos []*domain.Video) interface{} {
	// Not used in HTTP API - videos are accessed via gRPC
	return &HTTPResponse{
		StatusCode: http.StatusNotImplemented,
		Body: generated.Error{
			Code:    "NOT_IMPLEMENTED",
			Message: "Video operations are available via gRPC API",
		},
	}
}

// PresentCollectTrendingResult presents the result of collecting trending videos
func (p *videoPresenter) PresentCollectTrendingResult(result interface{}) interface{} {
	r, ok := result.(*input.CollectTrendingResult)
	if !ok {
		return p.PresentError(domain.ErrInvalidInput)
	}

	return &HTTPResponse{
		StatusCode: http.StatusOK,
		Body: generated.CollectTrendingResponse{
			VideosProcessed: int32(r.VideosProcessed),
			VideosAdded:     int32(r.VideosAdded),
			Duration:        r.Duration.String(),
		},
	}
}

// PresentCollectSubscriptionsResult presents the result of collecting subscription videos
func (p *videoPresenter) PresentCollectSubscriptionsResult(result interface{}) interface{} {
	r, ok := result.(*input.CollectSubscriptionsResult)
	if !ok {
		return p.PresentError(domain.ErrInvalidInput)
	}

	return &HTTPResponse{
		StatusCode: http.StatusOK,
		Body: generated.CollectSubscriptionsResponse{
			ChannelsProcessed: int32(r.ChannelsProcessed),
			VideosProcessed:   int32(r.VideosProcessed),
			VideosAdded:       int32(r.VideosAdded),
			Duration:          r.Duration.String(),
		},
	}
}

// PresentError presents an error
func (p *videoPresenter) PresentError(err error) interface{} {
	switch err {
	case domain.ErrVideoNotFound:
		return &HTTPResponse{
			StatusCode: http.StatusNotFound,
			Body: generated.Error{
				Code:    "VIDEO_NOT_FOUND",
				Message: err.Error(),
			},
		}
	default:
		return &HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body: generated.Error{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		}
	}
}