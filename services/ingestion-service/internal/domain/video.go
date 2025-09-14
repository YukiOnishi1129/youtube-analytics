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
	CreatedAt        time.Time
	UpdatedAt        *time.Time
	DeletedAt        *time.Time
	
	// Snapshots that need to be persisted (transient field)
	newSnapshots []*VideoSnapshot
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
		CreatedAt:        time.Now(),
	}, nil
}

// Update updates video metadata
func (v *Video) Update(title string) {
	v.Title = title
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

// AddSnapshot adds a new snapshot to the video
func (v *Video) AddSnapshot(snapshot *VideoSnapshot) error {
	// Validate snapshot belongs to this video
	if snapshot.VideoID != v.ID {
		return errors.New("snapshot does not belong to this video")
	}
	
	// Validate snapshot with video
	if err := snapshot.ValidateWithVideo(v); err != nil {
		return err
	}
	
	// Add to new snapshots list
	v.newSnapshots = append(v.newSnapshots, snapshot)
	return nil
}

// GetNewSnapshots returns snapshots that need to be persisted
func (v *Video) GetNewSnapshots() []*VideoSnapshot {
	return v.newSnapshots
}

// ClearNewSnapshots clears the new snapshots after persistence
func (v *Video) ClearNewSnapshots() {
	v.newSnapshots = nil
}