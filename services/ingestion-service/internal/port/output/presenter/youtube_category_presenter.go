package presenter

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
)

// YouTubeCategoryPresenter is the interface for presenting YouTube category data
type YouTubeCategoryPresenter interface {
	PresentYouTubeCategory(category *domain.YouTubeCategory) interface{}
	PresentYouTubeCategories(categories []*domain.YouTubeCategory) interface{}
	PresentYouTubeCategoryCreated(category *domain.YouTubeCategory) interface{}
	PresentYouTubeCategoryUpdated(category *domain.YouTubeCategory) interface{}
	PresentError(err error) interface{}
}