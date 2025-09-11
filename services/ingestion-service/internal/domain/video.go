package domain

import (
	"errors"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

var (
	ErrEmptyYouTubeVideoID     = errors.New("youtube video id cannot be empty")
	ErrVideoEmptyChannelID     = errors.New("youtube channel id cannot be empty")
	ErrInvalidVideoPublishedAt = errors.New("published at must be in the past")
)

// Video represents a YouTube video entity
type Video struct {
	ID               valueobject.UUID
	YouTubeVideoID   valueobject.YouTubeVideoID
	ChannelID        valueobject.UUID
	YouTubeChannelID valueobject.YouTubeChannelID
	Title            string
	PublishedAt      time.Time
	CategoryID       valueobject.CategoryID
	ThumbnailURL     string
	VideoURL         string
	CreatedAt        time.Time
	UpdatedAt        *time.Time
	DeletedAt        *time.Time
}

// NewVideo creates a new video
func NewVideo(
	id valueobject.UUID,
	youtubeVideoID valueobject.YouTubeVideoID,
	channelID valueobject.UUID,
	youtubeChannelID valueobject.YouTubeChannelID,
	title string,
	publishedAt time.Time,
	categoryID valueobject.CategoryID,
	thumbnailURL string,
	videoURL string,
) (*Video, error) {
	if youtubeVideoID == "" {
		return nil, ErrEmptyYouTubeVideoID
	}

	if youtubeChannelID == "" {
		return nil, ErrVideoEmptyChannelID
	}

	if publishedAt.After(time.Now()) {
		return nil, ErrInvalidVideoPublishedAt
	}

	return &Video{
		ID:               id,
		YouTubeVideoID:   youtubeVideoID,
		ChannelID:        channelID,
		YouTubeChannelID: youtubeChannelID,
		Title:            title,
		PublishedAt:      publishedAt,
		CategoryID:       categoryID,
		ThumbnailURL:     thumbnailURL,
		VideoURL:         videoURL,
		CreatedAt:        time.Now(),
	}, nil
}

// Update updates video metadata
func (v *Video) Update(title string, thumbnailURL string, videoURL string) {
	v.Title = title
	v.ThumbnailURL = thumbnailURL
	v.VideoURL = videoURL
	now := time.Now()
	v.UpdatedAt = &now
}

// Delete performs soft delete
func (v *Video) Delete() {
	now := time.Now()
	v.DeletedAt = &now
	v.UpdatedAt = &now
}

// IsDeleted checks if the video is deleted
func (v *Video) IsDeleted() bool {
	return v.DeletedAt != nil
}