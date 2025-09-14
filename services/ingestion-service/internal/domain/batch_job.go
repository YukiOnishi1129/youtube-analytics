package domain

import (
	"errors"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// JobType represents the type of batch job
type JobType string

const (
	JobTypeCollectTrending      JobType = "collect_trending"
	JobTypeRenewSubscriptions   JobType = "renew_subscriptions"
	JobTypeCollectSnapshots     JobType = "collect_snapshots"
)

// JobStatus represents the status of a batch job
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

var (
	ErrInvalidJobType   = errors.New("invalid job type")
	ErrInvalidJobStatus = errors.New("invalid job status")
)

// BatchJob represents a batch job execution record
type BatchJob struct {
	ID            valueobject.UUID
	JobType       JobType
	Status        JobStatus
	Parameters    map[string]interface{}
	StartedAt     *time.Time
	CompletedAt   *time.Time
	ErrorMessage  string
	Statistics    map[string]interface{}
	CreatedAt     time.Time
}

// NewBatchJob creates a new batch job
func NewBatchJob(
	id valueobject.UUID,
	jobType JobType,
	parameters map[string]interface{},
) (*BatchJob, error) {
	if !isValidJobType(jobType) {
		return nil, ErrInvalidJobType
	}

	return &BatchJob{
		ID:         id,
		JobType:    jobType,
		Status:     JobStatusPending,
		Parameters: parameters,
		CreatedAt:  time.Now(),
	}, nil
}

// Start marks the job as started
func (j *BatchJob) Start() error {
	if j.Status != JobStatusPending {
		return errors.New("job can only be started from pending status")
	}

	j.Status = JobStatusRunning
	now := time.Now()
	j.StartedAt = &now
	return nil
}

// Complete marks the job as completed with statistics
func (j *BatchJob) Complete(statistics map[string]interface{}) error {
	if j.Status != JobStatusRunning {
		return errors.New("job can only be completed from running status")
	}

	j.Status = JobStatusCompleted
	now := time.Now()
	j.CompletedAt = &now
	j.Statistics = statistics
	return nil
}

// Fail marks the job as failed with an error message
func (j *BatchJob) Fail(errorMessage string) error {
	if j.Status != JobStatusRunning {
		return errors.New("job can only be failed from running status")
	}

	j.Status = JobStatusFailed
	now := time.Now()
	j.CompletedAt = &now
	j.ErrorMessage = errorMessage
	return nil
}

func isValidJobType(jobType JobType) bool {
	switch jobType {
	case JobTypeCollectTrending, JobTypeRenewSubscriptions, JobTypeCollectSnapshots:
		return true
	default:
		return false
	}
}