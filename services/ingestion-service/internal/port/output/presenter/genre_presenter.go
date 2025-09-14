package presenter

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
)

// GenrePresenter is the interface for presenting genre data
type GenrePresenter interface {
	PresentGenre(genre *domain.Genre) interface{}
	PresentGenres(genres []*domain.Genre) interface{}
	PresentGenreCreated(genre *domain.Genre) interface{}
	PresentGenreUpdated(genre *domain.Genre) interface{}
	PresentGenreEnabled(genre *domain.Genre) interface{}
	PresentGenreDisabled(genre *domain.Genre) interface{}
	PresentError(err error) interface{}
}