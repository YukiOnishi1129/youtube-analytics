package cloudtasks

import (
	"context"
	"fmt"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// taskScheduler is a Google Cloud Tasks implementation of TaskScheduler
type taskScheduler struct {
	client     *cloudtasks.Client
	projectID  string
	location   string
	queueName  string
	serviceURL string
}

// NewTaskScheduler creates a new Cloud Tasks scheduler
func NewTaskScheduler(projectID, location, queueName, serviceURL string) (gateway.TaskScheduler, error) {
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud tasks client: %w", err)
	}

	return &taskScheduler{
		client:     client,
		projectID:  projectID,
		location:   location,
		queueName:  queueName,
		serviceURL: serviceURL,
	}, nil
}

// Schedule schedules a snapshot task
func (s *taskScheduler) Schedule(ctx context.Context, videoID valueobject.UUID, cp valueobject.CheckpointHour, eta time.Time) error {
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", s.projectID, s.location, s.queueName)

	req := &cloudtaskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &cloudtaskspb.Task{
			Name: fmt.Sprintf("%s/tasks/snapshot-%s-%d", queuePath, videoID, cp),
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					HttpMethod: cloudtaskspb.HttpMethod_POST,
					Url:        fmt.Sprintf("%s/internal/snapshots/%s/%d", s.serviceURL, videoID, cp),
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
			ScheduleTime: timestamppb.New(eta),
		},
	}

	_, err := s.client.CreateTask(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// Cancel cancels a scheduled task
func (s *taskScheduler) Cancel(ctx context.Context, videoID valueobject.UUID, cp valueobject.CheckpointHour) error {
	taskPath := fmt.Sprintf("projects/%s/locations/%s/queues/%s/tasks/snapshot-%s-%d",
		s.projectID, s.location, s.queueName, videoID, cp)

	req := &cloudtaskspb.DeleteTaskRequest{
		Name: taskPath,
	}

	err := s.client.DeleteTask(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// ScheduleSnapshot schedules a snapshot task
func (s *taskScheduler) ScheduleSnapshot(ctx context.Context, task *domain.SnapshotTask) error {
	// Calculate ETA based on checkpoint hour
	eta := time.Now().Add(time.Duration(task.CheckpointHour) * time.Hour)
	return s.Schedule(ctx, task.VideoID, task.CheckpointHour, eta)
}

// Close closes the cloud tasks client
func (s *taskScheduler) Close() error {
	return s.client.Close()
}
