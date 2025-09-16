package domain

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// YouTubeCategory represents YouTube category reference data
type YouTubeCategory struct {
	ID         valueobject.CategoryID // YouTube's category ID (e.g., 27, 28)
	Name       string                 // Category name (e.g., "Education", "Science & Technology")
	Assignable bool                   // Whether videos can be assigned to this category
}

// NewYouTubeCategory creates a new YouTubeCategory
func NewYouTubeCategory(id valueobject.CategoryID, name string, assignable bool) (*YouTubeCategory, error) {
	if name == "" {
		return nil, ErrInvalidInput
	}

	return &YouTubeCategory{
		ID:         id,
		Name:       name,
		Assignable: assignable,
	}, nil
}

// UpdateCategory updates category information
func (c *YouTubeCategory) UpdateCategory(name string, assignable bool) error {
	if name == "" {
		return ErrInvalidInput
	}

	c.Name = name
	c.Assignable = assignable
	return nil
}