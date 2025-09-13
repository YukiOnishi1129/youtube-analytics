package cloudtasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// cloudTasksGateway implements CloudTasksGateway interface
type cloudTasksGateway struct {
	client    *cloudtasks.Client
	projectID string
	location  string
	queueName string
	baseURL   string
}

// NewCloudTasksGateway creates a new Cloud Tasks gateway
func NewCloudTasksGateway(projectID, location, queueName, baseURL string) (gateway.CloudTasksGateway, error) {
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud tasks client: %w", err)
	}

	return &cloudTasksGateway{
		client:    client,
		projectID: projectID,
		location:  location,
		queueName: queueName,
		baseURL:   baseURL,
	}, nil
}

// CreateSnapshotTask creates a task to capture video snapshot
func (g *cloudTasksGateway) CreateSnapshotTask(ctx context.Context, videoID valueobject.UUID, checkpointHour valueobject.CheckpointHour, scheduleTime time.Time) error {
	// Build the queue path
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", g.projectID, g.location, g.queueName)

	// Create the task payload
	payload := map[string]interface{}{
		"videoId":        string(videoID),
		"checkpointHour": int(checkpointHour),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Build the task
	req := &cloudtaskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &cloudtaskspb.Task{
			Name: fmt.Sprintf("%s/tasks/snapshot-%s-%d", queuePath, videoID, checkpointHour),
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					HttpMethod: cloudtaskspb.HttpMethod_POST,
					Url:        fmt.Sprintf("%s/tasks/snapshot", g.baseURL),
					AuthorizationHeader: &cloudtaskspb.HttpRequest_OidcToken{
						OidcToken: &cloudtaskspb.OidcToken{
							ServiceAccountEmail: fmt.Sprintf("ingestion-service@%s.iam.gserviceaccount.com", g.projectID),
						},
					},
					Body: payloadBytes,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
			ScheduleTime: timestamppb.New(scheduleTime),
		},
	}

	// Create the task
	_, err = g.client.CreateTask(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// DeleteSnapshotTask deletes a scheduled snapshot task
func (g *cloudTasksGateway) DeleteSnapshotTask(ctx context.Context, videoID valueobject.UUID, checkpointHour valueobject.CheckpointHour) error {
	// Build the task path
	taskPath := fmt.Sprintf("projects/%s/locations/%s/queues/%s/tasks/snapshot-%s-%d",
		g.projectID, g.location, g.queueName, videoID, checkpointHour)

	req := &cloudtaskspb.DeleteTaskRequest{
		Name: taskPath,
	}

	err := g.client.DeleteTask(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// Close closes the client connection
func (g *cloudTasksGateway) Close() error {
	return g.client.Close()
}