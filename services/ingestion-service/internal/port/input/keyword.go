package input

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/google/uuid"
)

// KeywordInputPort is the interface for keyword use cases
type KeywordInputPort interface {
	ListKeywords(ctx context.Context) ([]*domain.Keyword, error)
	CreateKeyword(ctx context.Context, input *CreateKeywordInput) (*domain.Keyword, error)
	GetKeyword(ctx context.Context, keywordID uuid.UUID) (*domain.Keyword, error)
	UpdateKeyword(ctx context.Context, input *UpdateKeywordInput) (*domain.Keyword, error)
	EnableKeyword(ctx context.Context, keywordID uuid.UUID) (*domain.Keyword, error)
	DisableKeyword(ctx context.Context, keywordID uuid.UUID) (*domain.Keyword, error)
	DeleteKeyword(ctx context.Context, keywordID uuid.UUID) error
}

// CreateKeywordInput represents the input for creating a keyword
type CreateKeywordInput struct {
	Name        string
	FilterType  string
	Pattern     string
	Description *string
}

// UpdateKeywordInput represents the input for updating a keyword
type UpdateKeywordInput struct {
	KeywordID   uuid.UUID
	Name        string
	FilterType  string
	Pattern     string
	Enabled     bool
	Description *string
}
