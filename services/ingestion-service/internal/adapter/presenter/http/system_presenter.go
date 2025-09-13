package http

import (
	"net/http"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/driver/http/generated"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/presenter"
)

// systemPresenter implements presenter.SystemPresenter for HTTP
type systemPresenter struct{}

// NewSystemPresenter creates a new HTTP system presenter
func NewSystemPresenter() presenter.SystemPresenter {
	return &systemPresenter{}
}

// PresentSnapshot presents a single snapshot
// Note: This is for HTTP REST API, not used in current implementation
// as snapshot operations are exposed via gRPC only
func (p *systemPresenter) PresentSnapshot(snapshot *domain.VideoSnapshot) interface{} {
	// Not used in HTTP API - snapshots are accessed via gRPC
	return &HTTPResponse{
		StatusCode: http.StatusNotImplemented,
		Body: generated.Error{
			Code:    "NOT_IMPLEMENTED",
			Message: "Snapshot operations are available via gRPC API",
		},
	}
}

// PresentSnapshots presents multiple snapshots
// Note: This is for HTTP REST API, not used in current implementation
// as snapshot operations are exposed via gRPC only
func (p *systemPresenter) PresentSnapshots(snapshots []*domain.VideoSnapshot) interface{} {
	// Not used in HTTP API - snapshots are accessed via gRPC
	return &HTTPResponse{
		StatusCode: http.StatusNotImplemented,
		Body: generated.Error{
			Code:    "NOT_IMPLEMENTED",
			Message: "Snapshot operations are available via gRPC API",
		},
	}
}

// PresentScheduleSnapshotsResult presents the result of scheduling snapshots
func (p *systemPresenter) PresentScheduleSnapshotsResult(result interface{}) interface{} {
	r, ok := result.(*input.ScheduleSnapshotsResult)
	if !ok {
		return p.PresentError(domain.ErrInvalidInput)
	}

	return &HTTPResponse{
		StatusCode: http.StatusAccepted,
		Body: generated.ScheduleSnapshotsResponse{
			TasksScheduled:  int32(r.TasksScheduled),
			VideosProcessed: int32(r.VideosProcessed),
			Duration:        r.Duration.String(),
		},
	}
}

// PresentError presents an error
func (p *systemPresenter) PresentError(err error) interface{} {
	switch err {
	case domain.ErrSnapshotNotFound:
		return &HTTPResponse{
			StatusCode: http.StatusNotFound,
			Body: generated.Error{
				Code:    "SNAPSHOT_NOT_FOUND",
				Message: err.Error(),
			},
		}
	default:
		return &HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body: generated.Error{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		}
	}
}