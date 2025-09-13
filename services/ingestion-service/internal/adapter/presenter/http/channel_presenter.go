package http

import (
	"net/http"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/http/generated"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/presenter"
)

// channelPresenter implements presenter.ChannelPresenter for HTTP
type channelPresenter struct{}

// NewChannelPresenter creates a new HTTP channel presenter
func NewChannelPresenter() presenter.ChannelPresenter {
	return &channelPresenter{}
}

// PresentChannel presents a single channel
// Note: This is for HTTP REST API, not used in current implementation
// as channel operations are exposed via gRPC only
func (p *channelPresenter) PresentChannel(channel *domain.Channel) interface{} {
	// Not used in HTTP API - channels are accessed via gRPC
	return &HTTPResponse{
		StatusCode: http.StatusNotImplemented,
		Body: generated.Error{
			Code:    "NOT_IMPLEMENTED",
			Message: "Channel operations are available via gRPC API",
		},
	}
}

// PresentChannels presents multiple channels
// Note: This is for HTTP REST API, not used in current implementation
// as channel operations are exposed via gRPC only
func (p *channelPresenter) PresentChannels(channels []*domain.Channel) interface{} {
	// Not used in HTTP API - channels are accessed via gRPC
	return &HTTPResponse{
		StatusCode: http.StatusNotImplemented,
		Body: generated.Error{
			Code:    "NOT_IMPLEMENTED",
			Message: "Channel operations are available via gRPC API",
		},
	}
}

// PresentError presents an error
func (p *channelPresenter) PresentError(err error) interface{} {
	switch err {
	case domain.ErrChannelNotFound:
		return &HTTPResponse{
			StatusCode: http.StatusNotFound,
			Body: generated.Error{
				Code:    "CHANNEL_NOT_FOUND",
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

// PresentUpdateChannelsResult presents the result of updating channels
func (p *channelPresenter) PresentUpdateChannelsResult(result interface{}) interface{} {
	r, ok := result.(*input.UpdateChannelsResult)
	if !ok {
		return p.PresentError(domain.ErrInvalidInput)
	}

	return &HTTPResponse{
		StatusCode: http.StatusOK,
		Body: generated.UpdateChannelsResponse{
			ChannelsProcessed: int32(r.ChannelsProcessed),
			ChannelsUpdated:   int32(r.ChannelsUpdated),
			Duration:          r.Duration.String(),
		},
	}
}