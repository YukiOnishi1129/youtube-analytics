package presenter

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
)

// AuditLogPresenter is the interface for presenting audit log data
type AuditLogPresenter interface {
	PresentAuditLog(log *domain.AuditLog) interface{}
	PresentAuditLogs(logs []*domain.AuditLog) interface{}
	PresentAuditLogCreated(log *domain.AuditLog) interface{}
	PresentError(err error) interface{}
}