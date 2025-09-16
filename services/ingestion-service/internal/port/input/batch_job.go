package input

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/google/uuid"
)

// BatchJobInputPort is the interface for batch job use cases
type BatchJobInputPort interface {
	CreateBatchJob(ctx context.Context, input *CreateBatchJobInput) (*domain.BatchJob, error)
	StartBatchJob(ctx context.Context, jobID uuid.UUID) (*domain.BatchJob, error)
	CompleteBatchJob(ctx context.Context, jobID uuid.UUID, statistics map[string]interface{}) (*domain.BatchJob, error)
	FailBatchJob(ctx context.Context, jobID uuid.UUID, errorMessage string) (*domain.BatchJob, error)
	GetBatchJob(ctx context.Context, jobID uuid.UUID) (*domain.BatchJob, error)
	ListBatchJobsByType(ctx context.Context, jobType string, status string) ([]*domain.BatchJob, error)
	ListRecentBatchJobs(ctx context.Context, jobType *string, limit int) ([]*domain.BatchJob, error)
	GetRunningBatchJobs(ctx context.Context) ([]*domain.BatchJob, error)
}

// CreateBatchJobInput represents the input for creating a batch job
type CreateBatchJobInput struct {
	JobType    string
	Parameters map[string]interface{}
}