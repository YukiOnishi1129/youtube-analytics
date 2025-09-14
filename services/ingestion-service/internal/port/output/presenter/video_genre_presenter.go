package presenter

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
)

// VideoGenrePresenter is the interface for presenting video-genre association data
type VideoGenrePresenter interface {
	PresentVideoGenre(videoGenre *domain.VideoGenre) interface{}
	PresentVideoGenres(videoGenres []*domain.VideoGenre) interface{}
	PresentVideoGenreAssociated(videoGenre *domain.VideoGenre) interface{}
	PresentVideoGenresAssociated(videoGenres []*domain.VideoGenre) interface{}
	PresentVideoGenreDisassociated() interface{}
	PresentError(err error) interface{}
}