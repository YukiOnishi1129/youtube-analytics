package domain

import (
	"errors"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

var (
	ErrInvalidCheckpointHour = errors.New("invalid checkpoint hour")
	ErrInvalidMeasuredAt     = errors.New("measured at cannot be in the future")
	ErrInvalidSource         = errors.New("invalid source")
	ErrMeasuredAtBeforePublished = errors.New("measured at cannot be before video published at")
)

// VideoSnapshot represents a snapshot of video statistics at a checkpoint
type VideoSnapshot struct {
	ID                valueobject.UUID
	VideoID           valueobject.UUID
	CheckpointHour    valueobject.CheckpointHour
	MeasuredAt        time.Time
	ViewsCount        int64
	LikesCount        int64
	SubscriptionCount int64
	Source            valueobject.Source
	CreatedAt         time.Time
}

// SnapshotCounts holds the counts for a snapshot
type SnapshotCounts struct {
	ViewsCount        int64
	LikesCount        int64
	SubscriptionCount int64
}

// NewVideoSnapshot creates a new video snapshot
func NewVideoSnapshot(
	id valueobject.UUID,
	videoID valueobject.UUID,
	checkpointHour valueobject.CheckpointHour,
	measuredAt time.Time,
	counts SnapshotCounts,
	source valueobject.Source,
) (*VideoSnapshot, error) {
	if !checkpointHour.IsValid() {
		return nil, ErrInvalidCheckpointHour
	}

	if measuredAt.After(time.Now()) {
		return nil, ErrInvalidMeasuredAt
	}

	if !source.IsValid() {
		return nil, ErrInvalidSource
	}

	return &VideoSnapshot{
		ID:                id,
		VideoID:           videoID,
		CheckpointHour:    checkpointHour,
		MeasuredAt:        measuredAt,
		ViewsCount:        counts.ViewsCount,
		LikesCount:        counts.LikesCount,
		SubscriptionCount: counts.SubscriptionCount,
		Source:            source,
		CreatedAt:         time.Now(),
	}, nil
}

// ValidateWithVideo validates the snapshot against the video
func (s *VideoSnapshot) ValidateWithVideo(video *Video) error {
	if s.MeasuredAt.Before(video.PublishedAt) {
		return ErrMeasuredAtBeforePublished
	}
	return nil
}