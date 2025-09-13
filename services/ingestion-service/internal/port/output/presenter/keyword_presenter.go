package presenter

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
)

// KeywordPresenter is the output port for keyword presentation
type KeywordPresenter interface {
	PresentKeyword(keyword *domain.Keyword) interface{}
	PresentKeywords(keywords []*domain.Keyword) interface{}
	PresentDeleted() interface{}
	PresentError(err error) interface{}
}