package usecase

import (
	"context"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"github.com/google/uuid"
)

// batchJobUseCase implements the BatchJobInputPort interface
type batchJobUseCase struct {
	batchJobRepo gateway.BatchJobRepository
}

// NewBatchJobUseCase creates a new batch job use case
func NewBatchJobUseCase(batchJobRepo gateway.BatchJobRepository) input.BatchJobInputPort {
	return &batchJobUseCase{
		batchJobRepo: batchJobRepo,
	}
}

// CreateBatchJob creates a new batch job
func (u *batchJobUseCase) CreateBatchJob(ctx context.Context, input *input.CreateBatchJobInput) (*domain.BatchJob, error) {
	// Create domain object
	batchJob := &domain.BatchJob{
		ID:         valueobject.UUID(uuid.New().String()),
		JobType:    domain.JobType(input.JobType),
		Status:     domain.JobStatusPending,
		Parameters: input.Parameters,
	}

	// Save to repository
	if err := u.batchJobRepo.Save(ctx, batchJob); err != nil {
		return nil, err
	}

	return batchJob, nil
}

// StartBatchJob starts a batch job
func (u *batchJobUseCase) StartBatchJob(ctx context.Context, jobID uuid.UUID) (*domain.BatchJob, error) {
	// Find the job
	job, err := u.batchJobRepo.FindByID(ctx, valueobject.UUID(jobID.String()))
	if err != nil {
		return nil, err
	}

	// Update status
	job.Status = domain.JobStatusRunning
	now := time.Now()
	job.StartedAt = &now

	// Save to repository
	if err := u.batchJobRepo.Update(ctx, job); err != nil {
		return nil, err
	}

	return job, nil
}

// CompleteBatchJob completes a batch job
func (u *batchJobUseCase) CompleteBatchJob(ctx context.Context, jobID uuid.UUID, statistics map[string]interface{}) (*domain.BatchJob, error) {
	// Find the job
	job, err := u.batchJobRepo.FindByID(ctx, valueobject.UUID(jobID.String()))
	if err != nil {
		return nil, err
	}

	// Update status
	job.Status = domain.JobStatusCompleted
	now := time.Now()
	job.CompletedAt = &now
	job.Statistics = statistics

	// Save to repository
	if err := u.batchJobRepo.Update(ctx, job); err != nil {
		return nil, err
	}

	return job, nil
}

// FailBatchJob marks a batch job as failed
func (u *batchJobUseCase) FailBatchJob(ctx context.Context, jobID uuid.UUID, errorMessage string) (*domain.BatchJob, error) {
	// Find the job
	job, err := u.batchJobRepo.FindByID(ctx, valueobject.UUID(jobID.String()))
	if err != nil {
		return nil, err
	}

	// Update status
	job.Status = domain.JobStatusFailed
	now := time.Now()
	job.CompletedAt = &now
	job.ErrorMessage = errorMessage

	// Save to repository
	if err := u.batchJobRepo.Update(ctx, job); err != nil {
		return nil, err
	}

	return job, nil
}

// GetBatchJob gets a batch job by ID
func (u *batchJobUseCase) GetBatchJob(ctx context.Context, jobID uuid.UUID) (*domain.BatchJob, error) {
	return u.batchJobRepo.FindByID(ctx, valueobject.UUID(jobID.String()))
}

// ListBatchJobsByType lists batch jobs by type and status
func (u *batchJobUseCase) ListBatchJobsByType(ctx context.Context, jobType string, status string) ([]*domain.BatchJob, error) {
	return u.batchJobRepo.FindByTypeAndStatus(ctx, domain.JobType(jobType), domain.JobStatus(status))
}

// ListRecentBatchJobs lists recent batch jobs
func (u *batchJobUseCase) ListRecentBatchJobs(ctx context.Context, jobType *string, limit int) ([]*domain.BatchJob, error) {
	var jt *domain.JobType
	if jobType != nil {
		jobTypeVal := domain.JobType(*jobType)
		jt = &jobTypeVal
	}
	return u.batchJobRepo.FindRecent(ctx, jt, limit)
}

// GetRunningBatchJobs gets all running batch jobs
func (u *batchJobUseCase) GetRunningBatchJobs(ctx context.Context) ([]*domain.BatchJob, error) {
	return u.batchJobRepo.GetRunningJobs(ctx)
}