package service

import (
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
)

// SnapshotScheduler is a domain service for scheduling video snapshots
type SnapshotScheduler interface {
	ScheduleSnapshots(video *domain.Video) ([]ScheduledSnapshot, error)
}

// ScheduledSnapshot represents a scheduled snapshot task
type ScheduledSnapshot struct {
	VideoID        valueobject.UUID
	CheckpointHour valueobject.CheckpointHour
	ETA            time.Time
}

type snapshotScheduler struct{}

// NewSnapshotScheduler creates a new snapshot scheduler
func NewSnapshotScheduler() SnapshotScheduler {
	return &snapshotScheduler{}
}

// ScheduleSnapshots calculates the schedule for video snapshots
// Returns snapshots for +3h, +6h, +12h, +24h, +48h, +72h, +168h from published time
func (s *snapshotScheduler) ScheduleSnapshots(video *domain.Video) ([]ScheduledSnapshot, error) {
	var scheduled []ScheduledSnapshot
	now := time.Now()
	
	// Get all checkpoint hours after D0
	checkpointHours := valueobject.GetCheckpointHoursAfter(valueobject.CheckpointHour0)
	
	for _, cpHour := range checkpointHours {
		eta := video.PublishedAt.Add(time.Duration(cpHour) * time.Hour)
		
		// Skip if the ETA has already passed
		if eta.Before(now) {
			continue
		}
		
		scheduled = append(scheduled, ScheduledSnapshot{
			VideoID:        video.ID,
			CheckpointHour: cpHour,
			ETA:            eta,
		})
	}
	
	return scheduled, nil
}