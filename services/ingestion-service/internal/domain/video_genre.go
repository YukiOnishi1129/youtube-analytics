package domain

import (
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// VideoGenre represents the many-to-many relationship between videos and genres
type VideoGenre struct {
	ID        valueobject.UUID
	VideoID   valueobject.UUID
	GenreID   valueobject.UUID
	CreatedAt time.Time
}

// NewVideoGenre creates a new video-genre association
func NewVideoGenre(
	id valueobject.UUID,
	videoID valueobject.UUID,
	genreID valueobject.UUID,
) (*VideoGenre, error) {
	if videoID == "" {
		return nil, ErrInvalidInput
	}
	if genreID == "" {
		return nil, ErrInvalidInput
	}

	return &VideoGenre{
		ID:        id,
		VideoID:   videoID,
		GenreID:   genreID,
		CreatedAt: time.Now(),
	}, nil
}