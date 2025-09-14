package usecase

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/input"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
	"github.com/google/uuid"
)

// auditLogUseCase implements the AuditLogInputPort interface
type auditLogUseCase struct {
	auditLogRepo gateway.AuditLogRepository
}

// NewAuditLogUseCase creates a new audit log use case
func NewAuditLogUseCase(auditLogRepo gateway.AuditLogRepository) input.AuditLogInputPort {
	return &auditLogUseCase{
		auditLogRepo: auditLogRepo,
	}
}

// CreateAuditLog creates a new audit log entry
func (u *auditLogUseCase) CreateAuditLog(ctx context.Context, input *input.CreateAuditLogInput) (*domain.AuditLog, error) {
	// Create domain object
	auditLog := &domain.AuditLog{
		ID:           valueobject.UUID(uuid.New().String()),
		ActorID:      valueobject.UUID(input.ActorID.String()),
		ActorEmail:   input.ActorEmail,
		Action:       input.Action,
		ResourceType: input.ResourceType,
		ResourceID:   input.ResourceID,
		OldValues:    input.OldValues,
		NewValues:    input.NewValues,
		IPAddress:    input.IPAddress,
		UserAgent:    input.UserAgent,
	}

	// Save to repository
	if err := u.auditLogRepo.Save(ctx, auditLog); err != nil {
		return nil, err
	}

	return auditLog, nil
}

// ListAuditLogsByActor lists audit logs by actor ID
func (u *auditLogUseCase) ListAuditLogsByActor(ctx context.Context, actorID uuid.UUID, limit, offset int) ([]*domain.AuditLog, error) {
	return u.auditLogRepo.FindByActor(ctx, valueobject.UUID(actorID.String()), limit, offset)
}

// ListAuditLogsByResource lists audit logs by resource
func (u *auditLogUseCase) ListAuditLogsByResource(ctx context.Context, resourceType, resourceID string, limit, offset int) ([]*domain.AuditLog, error) {
	return u.auditLogRepo.FindByResource(ctx, resourceType, resourceID, limit, offset)
}

// ListRecentAuditLogs lists recent audit logs
func (u *auditLogUseCase) ListRecentAuditLogs(ctx context.Context, limit int) ([]*domain.AuditLog, error) {
	return u.auditLogRepo.FindRecent(ctx, limit)
}