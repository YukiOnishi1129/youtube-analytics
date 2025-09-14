package input

import (
	"context"
	"net"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/google/uuid"
)

// AuditLogInputPort is the interface for audit log use cases
type AuditLogInputPort interface {
	CreateAuditLog(ctx context.Context, input *CreateAuditLogInput) (*domain.AuditLog, error)
	ListAuditLogsByActor(ctx context.Context, actorID uuid.UUID, limit, offset int) ([]*domain.AuditLog, error)
	ListAuditLogsByResource(ctx context.Context, resourceType, resourceID string, limit, offset int) ([]*domain.AuditLog, error)
	ListRecentAuditLogs(ctx context.Context, limit int) ([]*domain.AuditLog, error)
}

// CreateAuditLogInput represents the input for creating an audit log
type CreateAuditLogInput struct {
	ActorID      uuid.UUID
	ActorEmail   string
	Action       string
	ResourceType string
	ResourceID   string
	OldValues    map[string]interface{}
	NewValues    map[string]interface{}
	IPAddress    net.IP
	UserAgent    string
}