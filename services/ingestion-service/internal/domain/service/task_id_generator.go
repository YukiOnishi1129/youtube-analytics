package service

import (
	"fmt"
	
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// TaskIDGenerator is a domain service for generating deterministic task IDs
type TaskIDGenerator interface {
	GenerateSnapshotTaskID(videoID valueobject.UUID, checkpointHour valueobject.CheckpointHour) string
}

type taskIDGenerator struct{}

// NewTaskIDGenerator creates a new task ID generator
func NewTaskIDGenerator() TaskIDGenerator {
	return &taskIDGenerator{}
}

// GenerateSnapshotTaskID generates a deterministic task ID for snapshot tasks
// Format: snap:{video_id}:{checkpoint_hour}
func (g *taskIDGenerator) GenerateSnapshotTaskID(videoID valueobject.UUID, checkpointHour valueobject.CheckpointHour) string {
	return fmt.Sprintf("snap:%s:%d", videoID, checkpointHour)
}