package mock

import (
	"context"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/port/output/gateway"
)

// keywordRepository is a mock implementation of KeywordRepository
type keywordRepository struct {
	keywords map[string]*domain.Keyword
}

// NewKeywordRepository creates a new mock keyword repository
func NewKeywordRepository() gateway.KeywordRepository {
	return &keywordRepository{
		keywords: make(map[string]*domain.Keyword),
	}
}

// Save saves a keyword
func (r *keywordRepository) Save(ctx context.Context, k *domain.Keyword) error {
	r.keywords[string(k.ID)] = k
	return nil
}

// FindAll finds all keywords
func (r *keywordRepository) FindAll(ctx context.Context, enabledOnly bool) ([]*domain.Keyword, error) {
	var result []*domain.Keyword
	for _, k := range r.keywords {
		if !enabledOnly || k.Enabled {
			result = append(result, k)
		}
	}
	return result, nil
}

// FindByID finds a keyword by ID
func (r *keywordRepository) FindByID(ctx context.Context, id valueobject.UUID) (*domain.Keyword, error) {
	k, exists := r.keywords[string(id)]
	if !exists {
		return nil, domain.ErrKeywordNotFound
	}
	return k, nil
}

// SoftDelete soft deletes a keyword
func (r *keywordRepository) SoftDelete(ctx context.Context, id valueobject.UUID) error {
	delete(r.keywords, string(id))
	return nil
}