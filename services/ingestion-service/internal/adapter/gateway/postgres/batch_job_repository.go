package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/adapter/gateway/postgres/sqlcgen"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// batchJobRepository implements gateway.BatchJobRepository interface
type batchJobRepository struct {
	*Repository
}

// NewBatchJobRepository creates a new batch job repository
func NewBatchJobRepository(repo *Repository) gateway.BatchJobRepository {
	return &batchJobRepository{Repository: repo}
}

// Save creates a new batch job
func (r *batchJobRepository) Save(ctx context.Context, job *domain.BatchJob) error {
	id, err := uuid.Parse(string(job.ID))
	if err != nil {
		return err
	}

	var parameters pqtype.NullRawMessage
	if job.Parameters != nil {
		data, err := json.Marshal(job.Parameters)
		if err != nil {
			return err
		}
		parameters = pqtype.NullRawMessage{RawMessage: data, Valid: true}
	}

	return r.q.CreateBatchJob(ctx, sqlcgen.CreateBatchJobParams{
		ID:         id,
		JobType:    string(job.JobType),
		Status:     string(job.Status),
		Parameters: parameters,
		CreatedAt:  job.CreatedAt,
	})
}

// Update updates an existing batch job
func (r *batchJobRepository) Update(ctx context.Context, job *domain.BatchJob) error {
	id, err := uuid.Parse(string(job.ID))
	if err != nil {
		return err
	}

	var statistics pqtype.NullRawMessage
	if job.Statistics != nil {
		data, err := json.Marshal(job.Statistics)
		if err != nil {
			return err
		}
		statistics = pqtype.NullRawMessage{RawMessage: data, Valid: true}
	}

	var startedAt, completedAt sql.NullTime
	if job.StartedAt != nil {
		startedAt = sql.NullTime{Time: *job.StartedAt, Valid: true}
	}
	if job.CompletedAt != nil {
		completedAt = sql.NullTime{Time: *job.CompletedAt, Valid: true}
	}

	return r.q.UpdateBatchJob(ctx, sqlcgen.UpdateBatchJobParams{
		ID:           id,
		Status:       string(job.Status),
		StartedAt:    startedAt,
		CompletedAt:  completedAt,
		ErrorMessage: sql.NullString{String: job.ErrorMessage, Valid: job.ErrorMessage != ""},
		Statistics:   statistics,
	})
}

// FindByID finds a batch job by ID
func (r *batchJobRepository) FindByID(ctx context.Context, id valueobject.UUID) (*domain.BatchJob, error) {
	jobID, err := uuid.Parse(string(id))
	if err != nil {
		return nil, err
	}

	row, err := r.q.GetBatchJobByID(ctx, jobID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainBatchJob(row), nil
}

// FindByTypeAndStatus finds batch jobs by type and status
func (r *batchJobRepository) FindByTypeAndStatus(ctx context.Context, jobType domain.JobType, status domain.JobStatus) ([]*domain.BatchJob, error) {
	rows, err := r.q.ListBatchJobsByTypeAndStatus(ctx, sqlcgen.ListBatchJobsByTypeAndStatusParams{
		JobType: string(jobType),
		Status:  string(status),
	})
	if err != nil {
		return nil, err
	}

	jobs := make([]*domain.BatchJob, len(rows))
	for i, row := range rows {
		jobs[i] = toDomainBatchJob(row)
	}

	return jobs, nil
}

// FindRecent finds recent batch jobs
func (r *batchJobRepository) FindRecent(ctx context.Context, jobType *domain.JobType, limit int) ([]*domain.BatchJob, error) {
	var jobTypeStr string
	if jobType != nil {
		jobTypeStr = string(*jobType)
	}

	rows, err := r.q.ListRecentBatchJobs(ctx, sqlcgen.ListRecentBatchJobsParams{
		Column1: jobTypeStr,
		Limit:   int32(limit),
	})
	if err != nil {
		return nil, err
	}

	jobs := make([]*domain.BatchJob, len(rows))
	for i, row := range rows {
		jobs[i] = toDomainBatchJob(row)
	}

	return jobs, nil
}

// GetRunningJobs finds all running batch jobs
func (r *batchJobRepository) GetRunningJobs(ctx context.Context) ([]*domain.BatchJob, error) {
	rows, err := r.q.ListRunningBatchJobs(ctx)
	if err != nil {
		return nil, err
	}

	jobs := make([]*domain.BatchJob, len(rows))
	for i, row := range rows {
		jobs[i] = toDomainBatchJob(row)
	}

	return jobs, nil
}

// toDomainBatchJob converts a database row to a domain batch job
func toDomainBatchJob(row sqlcgen.IngestionBatchJob) *domain.BatchJob {
	job := &domain.BatchJob{
		ID:        valueobject.UUID(row.ID.String()),
		JobType:   domain.JobType(row.JobType),
		Status:    domain.JobStatus(row.Status),
		CreatedAt: row.CreatedAt,
	}

	if row.Parameters.Valid {
		var parameters map[string]interface{}
		if err := json.Unmarshal(row.Parameters.RawMessage, &parameters); err == nil {
			job.Parameters = parameters
		}
	}

	if row.StartedAt.Valid {
		job.StartedAt = &row.StartedAt.Time
	}

	if row.CompletedAt.Valid {
		job.CompletedAt = &row.CompletedAt.Time
	}

	if row.ErrorMessage.Valid {
		job.ErrorMessage = row.ErrorMessage.String
	}

	if row.Statistics.Valid {
		var statistics map[string]interface{}
		if err := json.Unmarshal(row.Statistics.RawMessage, &statistics); err == nil {
			job.Statistics = statistics
		}
	}

	return job
}