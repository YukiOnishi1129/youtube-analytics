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
	ID                valueobject.UUID
	YouTubeChannelID  valueobject.YouTubeChannelID
	Title             string
	ThumbnailURL      string
	Description       string
	Country           string
	ViewCount         int64
	SubscriptionCount int64
	VideoCount        int64
	Subscribed        bool
	CreatedAt         time.Time
	UpdatedAt         *time.Time
	DeletedAt         *time.Time
	
	// Snapshots that need to be persisted (transient field)
	newSnapshots []*ChannelSnapshot
}

// NewChannel creates a new channel
func NewChannel(
	id valueobject.UUID,
	youtubeChannelID valueobject.YouTubeChannelID,
	title string,
	thumbnailURL string,
	description string,
	country string,
	viewCount int64,
	subscriptionCount int64,
	videoCount int64,
) (*Channel, error) {
	if youtubeChannelID == "" {
		return nil, ErrEmptyYouTubeChannelID
	}

	return &Channel{
		ID:                id,
		YouTubeChannelID:  youtubeChannelID,
		Title:             title,
		ThumbnailURL:      thumbnailURL,
		Description:       description,
		Country:           country,
		ViewCount:         viewCount,
		SubscriptionCount: subscriptionCount,
		VideoCount:        videoCount,
		Subscribed:        false,
		CreatedAt:         time.Now(),
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
func (c *Channel) UpdateProfile(title string, thumbnailURL string, description string, country string, viewCount int64, subscriptionCount int64, videoCount int64) {
	c.Title = title
	c.ThumbnailURL = thumbnailURL
	c.Description = description
	c.Country = country
	c.ViewCount = viewCount
	c.SubscriptionCount = subscriptionCount
	c.VideoCount = videoCount
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

// AddSnapshot adds a new snapshot to the channel
func (c *Channel) AddSnapshot(snapshot *ChannelSnapshot) error {
	// Validate snapshot belongs to this channel
	if snapshot.ChannelID != c.ID {
		return errors.New("snapshot does not belong to this channel")
	}
	
	// Add to new snapshots list
	c.newSnapshots = append(c.newSnapshots, snapshot)
	return nil
}

// GetNewSnapshots returns snapshots that need to be persisted
func (c *Channel) GetNewSnapshots() []*ChannelSnapshot {
	return c.newSnapshots
}

// ClearNewSnapshots clears the new snapshots after persistence
func (c *Channel) ClearNewSnapshots() {
	c.newSnapshots = nil
}

