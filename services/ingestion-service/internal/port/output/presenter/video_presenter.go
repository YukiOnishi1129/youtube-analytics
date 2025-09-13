package presenter

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
)

// VideoPresenter is the output port for video presentation
type VideoPresenter interface {
	PresentVideo(video *domain.Video) interface{}
	PresentVideos(videos []*domain.Video) interface{}
	PresentCollectTrendingResult(result interface{}) interface{}
	PresentCollectSubscriptionsResult(result interface{}) interface{}
	PresentError(err error) interface{}
}