package gateway

import "github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"

// UUIDGenerator generates unique identifiers
type UUIDGenerator interface {
	Generate() valueobject.UUID
}