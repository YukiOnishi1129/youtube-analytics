package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/YukiOnishi1129/youtube-analytics/services/ingestion-service/internal/domain/valueobject"
)

var (
	ErrEmptyName        = errors.New("keyword name cannot be empty")
	ErrEmptyPattern     = errors.New("keyword pattern cannot be empty")
	ErrInvalidFilterType = errors.New("invalid filter type")
)

// Keyword represents a filter keyword entity
type Keyword struct {
	ID          valueobject.UUID
	GenreID     valueobject.UUID       // Associated genre
	Name        string
	FilterType  valueobject.FilterType
	Pattern     string
	TargetField string                 // TITLE, DESCRIPTION, etc.
	Enabled     bool
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}

// NewKeyword creates a new keyword
func NewKeyword(
	id valueobject.UUID,
	genreID valueobject.UUID,
	name string,
	filterType valueobject.FilterType,
	pattern string,
	targetField string,
	description *string,
) (*Keyword, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyName
	}

	if strings.TrimSpace(pattern) == "" {
		return nil, ErrEmptyPattern
	}

	if !filterType.IsValid() {
		return nil, ErrInvalidFilterType
	}

	if targetField == "" {
		targetField = "TITLE" // Default to TITLE
	}

	return &Keyword{
		ID:          id,
		GenreID:     genreID,
		Name:        name,
		FilterType:  filterType,
		Pattern:     pattern,
		TargetField: targetField,
		Enabled:     true,
		Description: description,
		CreatedAt:   time.Now(),
	}, nil
}

// Update updates the keyword
func (k *Keyword) Update(name *string, filterType *valueobject.FilterType, pattern *string, description *string) error {
	if name != nil {
		if strings.TrimSpace(*name) == "" {
			return ErrEmptyName
		}
		k.Name = *name
	}

	if filterType != nil {
		if !filterType.IsValid() {
			return ErrInvalidFilterType
		}
		k.FilterType = *filterType
	}

	if pattern != nil {
		if strings.TrimSpace(*pattern) == "" {
			return ErrEmptyPattern
		}
		k.Pattern = *pattern
	}

	if description != nil {
		k.Description = description
	}

	now := time.Now()
	k.UpdatedAt = &now

	return nil
}

// Enable enables the keyword
func (k *Keyword) Enable() {
	k.Enabled = true
	now := time.Now()
	k.UpdatedAt = &now
}

// Disable disables the keyword
func (k *Keyword) Disable() {
	k.Enabled = false
	now := time.Now()
	k.UpdatedAt = &now
}

// Delete performs soft delete
func (k *Keyword) Delete() {
	now := time.Now()
	k.DeletedAt = &now
	k.UpdatedAt = &now
}

// IsDeleted checks if the keyword is deleted
func (k *Keyword) IsDeleted() bool {
	return k.DeletedAt != nil
}