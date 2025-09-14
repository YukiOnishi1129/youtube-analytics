package presenter

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
)

// BatchJobPresenter is the interface for presenting batch job data
type BatchJobPresenter interface {
	PresentBatchJob(job *domain.BatchJob) interface{}
	PresentBatchJobs(jobs []*domain.BatchJob) interface{}
	PresentBatchJobCreated(job *domain.BatchJob) interface{}
	PresentBatchJobStarted(job *domain.BatchJob) interface{}
	PresentBatchJobCompleted(job *domain.BatchJob) interface{}
	PresentBatchJobFailed(job *domain.BatchJob) interface{}
	PresentError(err error) interface{}
}