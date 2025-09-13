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
	Name        string
	FilterType  valueobject.FilterType
	Pattern     string
	Enabled     bool
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}

// NewKeyword creates a new keyword
func NewKeyword(
	id valueobject.UUID,
	name string,
	filterType valueobject.FilterType,
	pattern string,
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

	return &Keyword{
		ID:          id,
		Name:        name,
		FilterType:  filterType,
		Pattern:     pattern,
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