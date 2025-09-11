package input

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// KeywordUseCase defines the interface for keyword management use cases
type KeywordUseCase interface {
	ListKeywords(ctx context.Context) ([]*domain.Keyword, error)
	CreateKeyword(ctx context.Context, input CreateKeywordInput) (*domain.Keyword, error)
	UpdateKeyword(ctx context.Context, input UpdateKeywordInput) (*domain.Keyword, error)
	EnableKeyword(ctx context.Context, id valueobject.UUID) (*domain.Keyword, error)
	DisableKeyword(ctx context.Context, id valueobject.UUID) (*domain.Keyword, error)
	DeleteKeyword(ctx context.Context, id valueobject.UUID) error
}

// CreateKeywordInput represents the input for creating a keyword
type CreateKeywordInput struct {
	Name        string
	FilterType  valueobject.FilterType
	Description *string
}

// UpdateKeywordInput represents the input for updating a keyword
type UpdateKeywordInput struct {
	ID          valueobject.UUID
	Name        *string
	FilterType  *valueobject.FilterType
	Description *string
}