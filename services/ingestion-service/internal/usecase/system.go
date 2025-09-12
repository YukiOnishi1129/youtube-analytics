package usecase

import (
	"context"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/service"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

type systemUseCase struct {
	videoRepo        gateway.VideoRepository
	taskScheduler    gateway.TaskScheduler
	snapshotScheduler service.SnapshotScheduler
}

func NewSystemUseCase(
	videoRepo gateway.VideoRepository,
	taskScheduler gateway.TaskScheduler,
	snapshotScheduler service.SnapshotScheduler,
) input.SystemInputPort {
	return &systemUseCase{
		videoRepo:         videoRepo,
		taskScheduler:     taskScheduler,
		snapshotScheduler: snapshotScheduler,
	}
}

func (u *systemUseCase) ScheduleSnapshots(ctx context.Context) (*input.ScheduleSnapshotsResult, error) {
	// Get active videos (videos that need snapshots)
	activeVideos, err := u.videoRepo.ListActive(ctx, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, err
	}

	tasksScheduled := 0
	for _, video := range activeVideos {
		// Determine checkpoint hours for this video
		checkpoints := u.snapshotScheduler.DetermineCheckpoints(video)
		
		for _, checkpoint := range checkpoints {
			// Schedule snapshot task
			task := &domain.SnapshotTask{
				VideoID:        video.ID,
				CheckpointHour: checkpoint,
				ScheduledAt:    time.Now(),
			}

			if err := u.taskScheduler.ScheduleSnapshot(ctx, task); err != nil {
				// Log error but continue with other tasks
				continue
			}
			tasksScheduled++
		}
	}

	return &input.ScheduleSnapshotsResult{
		VideosProcessed: len(activeVideos),
		TasksScheduled:  tasksScheduled,
	}, nil
}