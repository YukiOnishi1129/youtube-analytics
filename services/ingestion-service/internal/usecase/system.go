package usecase

import (
	"context"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/service"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"github.com/google/uuid"
)

type systemUseCase struct {
	videoRepo         gateway.VideoRepository
	snapshotRepo      gateway.VideoSnapshotRepository
	taskScheduler     gateway.TaskScheduler
	snapshotScheduler service.SnapshotScheduler
	youtubeAPI        gateway.YouTubeClient
}

func NewSystemUseCase(
	videoRepo gateway.VideoRepository,
	snapshotRepo gateway.VideoSnapshotRepository,
	taskScheduler gateway.TaskScheduler,
	snapshotScheduler service.SnapshotScheduler,
	youtubeAPI gateway.YouTubeClient,
) input.SystemInputPort {
	return &systemUseCase{
		videoRepo:         videoRepo,
		snapshotRepo:      snapshotRepo,
		taskScheduler:     taskScheduler,
		snapshotScheduler: snapshotScheduler,
		youtubeAPI:        youtubeAPI,
	}
}

func (u *systemUseCase) ScheduleSnapshots(ctx context.Context) (*input.ScheduleSnapshotsResult, error) {
	start := time.Now()

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
		Duration:        time.Since(start),
	}, nil
}

func (u *systemUseCase) CreateSnapshot(ctx context.Context, input *input.CreateSnapshotInput) (*domain.VideoSnapshot, error) {
	// Get video
	video, err := u.videoRepo.GetByID(ctx, valueobject.UUID(input.VideoID.String()))
	if err != nil {
		return nil, err
	}

	// Check if snapshot already exists
	exists, err := u.snapshotRepo.Exists(ctx, video.ID, valueobject.CheckpointHour(input.CheckpointHour))
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrSnapshotAlreadyExists
	}

	// Fetch current stats from YouTube API
	stats, err := u.youtubeAPI.GetVideoStatistics(ctx, string(video.YouTubeVideoID))
	if err != nil {
		return nil, err
	}

	// Create snapshot
	snapshot, err := domain.NewVideoSnapshot(
		valueobject.UUID(uuid.New().String()),
		video.ID,
		valueobject.CheckpointHour(input.CheckpointHour),
		time.Now(),
		domain.SnapshotCounts{
			ViewsCount:        stats.ViewCount,
			LikesCount:        stats.LikeCount,
			SubscriptionCount: stats.CommentCount,
		},
		valueobject.Source("youtube_api"),
	)
	if err != nil {
		return nil, err
	}

	// Add snapshot to video aggregate
	if err := video.AddSnapshot(snapshot); err != nil {
		return nil, err
	}

	// Save video with snapshots
	if err := u.videoRepo.SaveWithSnapshots(ctx, video); err != nil {
		return nil, err
	}

	return snapshot, nil
}

func (u *systemUseCase) GetVideoSnapshots(ctx context.Context, videoID uuid.UUID) ([]*domain.VideoSnapshot, error) {
	// Verify video exists
	_, err := u.videoRepo.GetByID(ctx, valueobject.UUID(videoID.String()))
	if err != nil {
		return nil, err
	}

	// Get all snapshots for the video
	return u.snapshotRepo.ListByVideoID(ctx, valueobject.UUID(videoID.String()))
}
