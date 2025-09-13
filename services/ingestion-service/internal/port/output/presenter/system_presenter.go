package presenter

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
)

// SystemPresenter is the output port for system presentation
type SystemPresenter interface {
	PresentSnapshot(snapshot *domain.VideoSnapshot) interface{}
	PresentSnapshots(snapshots []*domain.VideoSnapshot) interface{}
	PresentScheduleSnapshotsResult(result interface{}) interface{}
	PresentError(err error) interface{}
}