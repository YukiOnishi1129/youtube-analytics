package uuid

import (
	"github.com/google/uuid"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

// generator implements gateway.UUIDGenerator
type generator struct{}

// NewGenerator creates a new UUID generator
func NewGenerator() gateway.UUIDGenerator {
	return &generator{}
}

// Generate generates a new UUID
func (g *generator) Generate() valueobject.UUID {
	return valueobject.UUID(uuid.New().String())
}