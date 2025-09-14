package domain

import (
	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

// Genre represents a collection target with region, categories, and language settings
type Genre struct {
	ID          valueobject.UUID
	Code        string                   // e.g., "engineering_jp", "engineering_en"
	Name        string                   // e.g., "Engineering (JP)"
	Language    string                   // e.g., "ja", "en"
	RegionCode  string                   // e.g., "JP", "US"
	CategoryIDs []valueobject.CategoryID // e.g., [27, 28] for Education & Science
	Enabled     bool
}

// NewGenre creates a new Genre
func NewGenre(
	id valueobject.UUID,
	code, name, language, regionCode string,
	categoryIDs []valueobject.CategoryID,
) (*Genre, error) {
	if code == "" {
		return nil, ErrInvalidInput
	}
	if name == "" {
		return nil, ErrInvalidInput
	}
	if language == "" {
		return nil, ErrInvalidInput
	}
	if regionCode == "" {
		return nil, ErrInvalidInput
	}
	if len(categoryIDs) == 0 {
		return nil, ErrInvalidInput
	}

	return &Genre{
		ID:          id,
		Code:        code,
		Name:        name,
		Language:    language,
		RegionCode:  regionCode,
		CategoryIDs: categoryIDs,
		Enabled:     true,
	}, nil
}

// Enable enables the genre for collection
func (g *Genre) Enable() error {
	g.Enabled = true
	return nil
}

// Disable disables the genre for collection
func (g *Genre) Disable() error {
	g.Enabled = false
	return nil
}

// Update updates the genre properties
func (g *Genre) Update(name string, categoryIDs []valueobject.CategoryID) error {
	if name == "" {
		return ErrInvalidInput
	}
	if len(categoryIDs) == 0 {
		return ErrInvalidInput
	}

	g.Name = name
	g.CategoryIDs = categoryIDs
	return nil
}