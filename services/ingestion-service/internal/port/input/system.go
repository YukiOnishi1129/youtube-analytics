package input

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// SystemUseCase defines the interface for system operations use cases
type SystemUseCase interface {
	CollectTrending(ctx context.Context, input CollectTrendingInput) (*CollectTrendingOutput, error)
	HealthCheck(ctx context.Context, input HealthCheckInput) (*HealthCheckOutput, error)
	Warm(ctx context.Context) error
}

// CollectTrendingInput represents the input for collecting trending videos
type CollectTrendingInput struct {
	Region      string
	CategoryIDs []valueobject.CategoryID
	Pages       int
}

// CollectTrendingOutput represents the output for collecting trending videos
type CollectTrendingOutput struct {
	Collected int
	Adopted   int
}

// HealthCheckInput represents the input for WebSub health check
type HealthCheckInput struct {
	Mode         string // "subscribe" or "unsubscribe"
	Topic        string
	Challenge    string
	LeaseSeconds int
}

// HealthCheckOutput represents the output for WebSub health check
type HealthCheckOutput struct {
	Challenge string
}