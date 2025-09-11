package domain

import (
	"errors"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

var (
	ErrEmptyYouTubeChannelID = errors.New("youtube channel id cannot be empty")
)

// Channel represents a YouTube channel entity
type Channel struct {
	ID               valueobject.UUID
	YouTubeChannelID valueobject.YouTubeChannelID
	Title            string
	ThumbnailURL     string
	Subscribed       bool
	CreatedAt        time.Time
	UpdatedAt        *time.Time
	DeletedAt        *time.Time
}

// NewChannel creates a new channel
func NewChannel(
	id valueobject.UUID,
	youtubeChannelID valueobject.YouTubeChannelID,
	title string,
	thumbnailURL string,
) (*Channel, error) {
	if youtubeChannelID == "" {
		return nil, ErrEmptyYouTubeChannelID
	}

	return &Channel{
		ID:               id,
		YouTubeChannelID: youtubeChannelID,
		Title:            title,
		ThumbnailURL:     thumbnailURL,
		Subscribed:       false,
		CreatedAt:        time.Now(),
	}, nil
}

// Subscribe subscribes to the channel
func (c *Channel) Subscribe() {
	c.Subscribed = true
	now := time.Now()
	c.UpdatedAt = &now
}

// Unsubscribe unsubscribes from the channel
func (c *Channel) Unsubscribe() {
	c.Subscribed = false
	now := time.Now()
	c.UpdatedAt = &now
}

// UpdateProfile updates channel profile information
func (c *Channel) UpdateProfile(title string, thumbnailURL string) {
	c.Title = title
	c.ThumbnailURL = thumbnailURL
	now := time.Now()
	c.UpdatedAt = &now
}

// Delete performs soft delete
func (c *Channel) Delete() {
	now := time.Now()
	c.DeletedAt = &now
	c.UpdatedAt = &now
}

// IsDeleted checks if the channel is deleted
func (c *Channel) IsDeleted() bool {
	return c.DeletedAt != nil
}

