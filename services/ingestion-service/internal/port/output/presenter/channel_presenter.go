package presenter

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
)

// ChannelPresenter is the output port for channel presentation
type ChannelPresenter interface {
	PresentChannel(channel *domain.Channel) interface{}
	PresentChannels(channels []*domain.Channel) interface{}
	PresentError(err error) interface{}
	PresentUpdateChannelsResult(result interface{}) interface{}
}